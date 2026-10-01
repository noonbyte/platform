.PHONY: release

release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release VERSION=v1.2.3"; \
		exit 1; \
	fi
	@if git rev-parse "$(VERSION)" >/dev/null 2>&1; then \
		echo "Error: tag $(VERSION) already exists"; \
		exit 1; \
	fi
	@echo "Releasing $(VERSION)..."
	git add .
	git commit -m "release: $(VERSION)"
	git tag "$(VERSION)"
	git push origin main
	git push origin "$(VERSION)"
	@echo "Released $(VERSION)"
