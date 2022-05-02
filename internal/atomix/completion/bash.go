// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package completion

import (
	"github.com/spf13/cobra"
	"io"
)

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	return cmd.GenBashCompletion(out)
}

const BashCompletionFunction = `
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
