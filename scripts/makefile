TEST_IMAGE = test-git-backups:test
SSH_AUTH_PATH = /tmp/ssh-auth.sock

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
test-image-ssh-key: build-image
	docker run --rm \
		-v ../test/test_image_backup/config_ssh_key.yaml:/app/config.yaml:ro \
		-v ../test/test_image_backup/ssh_key:/app/ssh_key:ro \
		-v ../test/test_image_backup/backup:/app/backup:rw \
		--name test-git-backups $(TEST_IMAGE)
	$(MAKE) remove-image

test-image-ssh-agent: build-image
	docker run --rm \
		--mount type=bind,src=../test/test_image_backup/known_hosts,dst=/home/appuser/.ssh/known_hosts,readonly \
		--mount type=bind,src=../test/test_image_backup/config_ssh_agent.yaml,dst=/app/config.yaml,readonly \
		--mount type=bind,src=../test/test_image_backup/backup,dst=/app/backup \
		--volume type=bind,src=$$SSH_AUTH_SOCK,dst=$(SSH_AUTH_PATH),readonly \
		--env SSH_AUTH_SOCK=$(SSH_AUTH_PATH) \
		--name test-git-backups $(TEST_IMAGE)
	$(MAKE) remove-image

build-image:
	docker build --tag $(TEST_IMAGE) --no-cache --file ../build/release/dockerfile ..

remove-image:
	docker image rm $(TEST_IMAGE)