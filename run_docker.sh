docker run -it --network=host -v $(pwd):/root/app --entrypoint "/run_server.sh" go-translate-server $@
