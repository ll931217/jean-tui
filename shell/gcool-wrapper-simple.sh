#!/usr/bin/env bash
# Simple shell wrapper for gcool
# Source this file in your shell rc file

gcool() {
    # Create a temp file for communication
    local temp_file=$(mktemp)

    # Set environment variable so gcool knows to write to file
    GCOOL_SWITCH_FILE="$temp_file" command gcool "$@"
    local exit_code=$?

    # Check if switch info was written
    if [ -f "$temp_file" ] && [ -s "$temp_file" ]; then
        # Read the switch info: path|branch|auto-claude
        local switch_info=$(cat "$temp_file")
        # Only remove if it's in /tmp (safety check)
        if [[ "$temp_file" == /tmp/* ]] || [[ "$temp_file" == /var/folders/* ]]; then
            rm "$temp_file"
        fi

        # Parse the info (using worktree_path instead of path to avoid PATH conflict)
        IFS='|' read -r worktree_path branch auto_claude <<< "$switch_info"

        # Check if we got valid data (has both pipes)
        if [[ "$switch_info" == *"|"*"|"* ]]; then
            # Sanitize branch name for session
            local session_name="gcool-${branch//[^a-zA-Z0-9\-_]/-}"

            echo "Switching to: $worktree_path"
            echo "Branch: $branch"
            echo "Session: $session_name"

            # Check if tmux is available
            if ! command -v tmux >/dev/null 2>&1; then
                cd "$worktree_path" || return
                return
            fi

            # Check if already in tmux
            if [ -n "$TMUX" ]; then
                cd "$worktree_path" || return
                echo "Already in tmux, just changed directory"
                return
            fi

            # Check if session exists
            if tmux has-session -t "$session_name" 2>/dev/null; then
                echo "Attaching to existing session: $session_name"
                exec tmux attach-session -t "$session_name"
            else
                echo "Creating new session: $session_name"
                # Try to start with Claude if requested
                if [ "$auto_claude" = "true" ] && command -v claude >/dev/null 2>&1; then
                    echo "Starting with Claude CLI..."
                    exec tmux new-session -s "$session_name" -c "$worktree_path" claude
                else
                    echo "Starting shell..."
                    exec tmux new-session -s "$session_name" -c "$worktree_path"
                fi
            fi
        fi
    else
        # No switch file, just clean up
        # Only remove if it's in /tmp (safety check)
        if [[ "$temp_file" == /tmp/* ]] || [[ "$temp_file" == /var/folders/* ]]; then
            rm -f "$temp_file"
        fi
    fi

    return $exit_code
}
