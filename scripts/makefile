generate:
	go generate ../...

test:
	go tool ginkgo -r ..

tidy:
	go vet ../...
	go tool staticcheck ../...

# The `../test/test_image_backup` folder must contain:
# * The private SSH key (`ssh_key`).
# * The config file (`config.yml`).
#
# The backup will be created in the `backup` subfolder.
test_image:
	docker build --tag test-git-backups:test --no-cache --file ../build/release/dockerfile ..
	docker run --rm \
		-v ../test/test_image_backup/config.yaml:/app/config.yaml:ro \
		-v ../test/test_image_backup/ssh_key:/app/ssh_key:ro \
		-v ../test/test_image_backup/backup:/app/backup:rw \
		--name test-git-backups test-git-backups:test
	docker image rm test-git-backups:test
