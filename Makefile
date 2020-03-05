PLUGINS := auto_attachment_action echo keyword_reaction movie_announcement ping random_direct_message random_reaction run_command
OUTPUT := ./build

clean:
	rm ${OUTPUT}/* || true

build-all:
	@for p in ${PLUGINS}; do \
		go build -buildmode=plugin -o ${OUTPUT}/$$p.so $$p/*.go ; \
	done
