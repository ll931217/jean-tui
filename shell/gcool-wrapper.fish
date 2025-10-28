#!/usr/bin/env fish
# Shell wrapper for gcool to enable directory switching and tmux session management (Fish shell)
# Source this file in your config.fish

function gcool
    # Create a temp file for communication
    set temp_file (mktemp)

    # Set environment variable so gcool knows to write to file
    set -x GCOOL_SWITCH_FILE $temp_file
    command gcool $argv
    set exit_code $status

    # Check if switch info was written
    if test -f "$temp_file" -a -s "$temp_file"
        # Read the switch info: path|branch|auto-claude
        set switch_info (cat $temp_file)
        rm $temp_file

        # Parse the info (using worktree_path instead of path to avoid PATH conflict)
        set parts (string split '|' $switch_info)

        # Check if we got valid data (has 3 parts)
        if test (count $parts) -eq 3
            set worktree_path $parts[1]
            set branch $parts[2]
            set auto_claude $parts[3]

            # Check if tmux is available
            if not command -v tmux &> /dev/null
                # No tmux, just cd
                cd $worktree_path
                echo "Switched to worktree: $branch"
                return
            end

            # Sanitize branch name for tmux session
            set session_name "gcool-"(string replace -ra '[^a-zA-Z0-9\-_]' '-' $branch)
            set session_name (string replace -ra '--+' '-' $session_name)
            set session_name (string trim -c '-' $session_name)

            # Check if already in a tmux session
            if test -n "$TMUX"
                # Already in tmux, just cd
                cd $worktree_path
                echo "Switched to worktree: $branch"
                echo "Note: Already in tmux. Session: $session_name would be available outside tmux."
                return
            end

            # Check if session exists
            if tmux has-session -t "$session_name" 2>/dev/null
                # Attach to existing session
                exec tmux attach-session -t "$session_name"
            else
                # Create new session
                if test "$auto_claude" = "true"
                    # Check if claude is available
                    if command -v claude &> /dev/null
                        # Start with claude
                        exec tmux new-session -s "$session_name" -c "$worktree_path" claude
                    else
                        # Fallback: start with shell and show message
                        exec tmux new-session -s "$session_name" -c "$worktree_path"
                    end
                else
                    # Start with shell
                    exec tmux new-session -s "$session_name" -c "$worktree_path"
                end
            end
        end
    else
        # No switch file, just clean up
        rm -f $temp_file
    end

    return $exit_code
end
