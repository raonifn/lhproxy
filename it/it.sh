#!/bin/bash -e

cmd_cleanup() {
  docker rm -f lhproxy_it_squid || true
  docker rm -f lhproxy_it_pipe || true
}

cmd_build() {
  lhproxy_it_target="${1:-final}"
  cd ..
  docker build -t lhproxy/it:dev --target "$lhproxy_it_target" -f it/Dockerfile .
  cd -
}

cmd_test() {
  trap cmd_cleanup EXIT
  cmd_cleanup
  docker run -d --rm --label lhproxy_dev --name lhproxy_it_squid \
    -p 3128:3128 -h lhproxy_it_squid lhproxy/it:dev /root/entrypoint/server.sh

  docker run -it --rm --label lhproxy_dev --name lhproxy_it_pipe \
    --network host lhproxy/it:dev /root/entrypoint/pipe.sh

  echo "SUCCESS"
}

cmd_it() {
  cmd_build
  cmd_test
}

cd "$(dirname "$0")"; _cmd="${1?"cmd is required"}"; shift; "cmd_${_cmd}" "$@"
