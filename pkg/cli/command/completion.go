// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"io"
	"os"
)

const bashCompletion = `

__atomix_override_flag_list=(--scope -s --database -d)
__atomix_override_flags()
{
    local ${__atomix_override_flag_list[*]##*-} two_word_of of var
    for w in "${words[@]}"; do
        if [ -n "${two_word_of}" ]; then
            eval "${two_word_of##*-}=\"${two_word_of}=\${w}\""
            two_word_of=
            continue
        fi
        for of in "${__atomix_override_flag_list[@]}"; do
            case "${w}" in
                ${of}=*)
                    eval "${of##*-}=\"${w}\""
                    ;;
                ${of})
                    two_word_of="${of}"
                    ;;
            esac
        done
    done
    for var in "${__atomix_override_flag_list[@]##*-}"; do
        if eval "test -n \"\$${var}\""; then
            eval "echo -n \${${var}}' '"
        fi
    done
}

__atomix_get_databases() {
    local atomix_output out
    if atomix_output=$(atomix get databases --no-headers 2>/dev/null); then
        out=($(echo "${atomix_output}" | awk '{print $2}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__atomix_get_scopes() {
    local atomix_output out
    if atomix_output=$(atomix get primitives --no-headers 2>/dev/null); then
        out=($(echo "${atomix_output}" | awk '{print $2}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__atomix_primitive_types() {
    echo "counter"
    echo "election"
    echo "indexed-map"
    echo "leader-latch"
    echo "list"
    echo "lock"
    echo "log"
    echo "map"
    echo "set"
    echo "value"
}

__atomix_get_primitive_types() {
    local atomix_output out
    if out=$(__atomix_primitive_types); then
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__atomix_get_primitives() {
    local atomix_output out
    if atomix_output=$(atomix get primitives $(__atomix_override_flags) --type=$1 --no-headers 2>/dev/null); then
        out=($(echo "${atomix_output}" | awk '{print $1}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__atomix_get_counters() {
    __atomix_get_primitives "counter"
}

__atomix_get_elections() {
    __atomix_get_primitives "election"
}

__atomix_get_indexed_maps() {
    __atomix_get_primitives "indexed-map"
}

__atomix_get_leader_latches() {
    __atomix_get_primitives "leader-latch"
}

__atomix_get_lists() {
    __atomix_get_primitives "list"
}

__atomix_get_locks() {
    __atomix_get_primitives "lock"
}

__atomix_get_logs() {
    __atomix_get_primitives "log"
}

__atomix_get_maps() {
    __atomix_get_primitives "map"
}

__atomix_get_sets() {
    __atomix_get_primitives "set"
}

__atomix_get_values() {
    __atomix_get_primitives "value"
}

__atomix_custom_func() {
    case ${last_command} in
        atomix_create_counter | atomix_delete_counter)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_counters
            fi
            return
            ;;
        atomix_create_election | atomix_delete_election)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_elections
            fi
            return
            ;;
        atomix_create_indexed_map | atomix_delete_indexed_map)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_indexed_maps
            fi
            return
            ;;
        atomix_create_leader_latch | atomix_delete_leader_latch)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_leader_latches
            fi
            return
            ;;
        atomix_create_list | atomix_delete_list)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_lists
            fi
            return
            ;;
        atomix_create_lock | atomix_delete_lock)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_locks
            fi
            return
            ;;
        atomix_create_log | atomix_delete_log)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_logs
            fi
            return
            ;;
        atomix_create_map | atomix_delete_map)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_maps
            fi
            return
            ;;
        atomix_create_set | atomix_delete_set)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_sets
            fi
            return
            ;;
        atomix_create_value | atomix_delete_value)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __atomix_get_values
            fi
            return
            ;;
        *)
            ;;
    esac
}
`

func newCompletionCommand() *cobra.Command {
	return &cobra.Command{
		Use:       "completion <shell>",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh"},
		Run:       runCompletionCommand,
	}
}

func runCompletionCommand(cmd *cobra.Command, args []string) {
	if args[0] == "bash" {
		if err := runCompletionBash(os.Stdout, cmd.Parent()); err != nil {
			ExitWithError(ExitError, err)
		}
	} else if args[0] == "zsh" {
		if err := runCompletionZsh(os.Stdout, cmd.Parent()); err != nil {
			ExitWithError(ExitError, err)
		}
	} else {
		ExitWithError(ExitError, errors.New("unsupported shell type "+args[0]))
	}
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	return cmd.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, cmd *cobra.Command) error {
	zsh_head := "#compdef atomix\n"

	out.Write([]byte(zsh_head))

	zsh_initialization := `
__atomix_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand

	source "$@"
}

__atomix_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift

		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__atomix_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}

__atomix_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?

	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}

__atomix_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}

__atomix_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}

__atomix_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}

__atomix_filedir() {
	local RET OLD_IFS w qw

	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi

	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"

	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"

	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__atomix_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}

__atomix_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
    	printf %q "$1"
    fi
}

autoload -U +X bashcompinit && bashcompinit

# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi

__atomix_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__atomix_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__atomix_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__atomix_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__atomix_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__atomix_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__atomix_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	out.Write([]byte(zsh_initialization))

	buf := new(bytes.Buffer)
	cmd.GenBashCompletion(buf)
	out.Write(buf.Bytes())

	zsh_tail := `
BASH_COMPLETION_EOF
}

__atomix_bash_source <(__atomix_convert_bash_to_zsh)
_complete atomix 2>/dev/null
`
	out.Write([]byte(zsh_tail))
	return nil
}
