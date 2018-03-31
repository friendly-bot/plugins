PLUGINS := keyword_reaction random_reaction use_invite_not_at random_direct_message auto_attachment_action

build-all:
	@for p in $(PLUGINS); do \
		go build -buildmode=plugin -o $$p/$$p.so $$p/*.go ; \
	done

clean-all:
	@for p in $(PLUGINS); do \
    		rm $$p/$$p.so; \
    	done

list-all:
	@for p in $(PLUGINS); do \
		echo "$$p: `pwd`/$$p/$$p.so"; \
	done
