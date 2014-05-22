# Functions that run specs


# Makes a full spec run with all the specs
function run_all_specs {
  for spec_file in `find $git_town_root_dir/spec -name *_spec.sh`; do
    local spec_name=${spec_file/$git_town_root_dir\//}
    echo_header "RUNNING $spec_name"
    run_spec $spec_file
  done
}


# Makes a full spec run with only the given spec
function run_single_spec {
  reset_test_repo
  run_spec $1
}


# Runs the given spec
#
# This method
function run_spec {
  main_branch_created=false
  remote_main_branch_created=false
  source $1
  unset -f before
  echo
  echo $underline"Cleanup up after '$current_SUT' specs"$nounderline
  run_after_function
  reset_test_repo
  echo $underline"Cleanup done"$nounderline
}
