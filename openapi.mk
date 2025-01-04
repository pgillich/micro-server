
openapi-server-ogen:
	docker run ${DOCKER_RUN_FLAGS} \
		${DOCKER_OGEN_PATH}/${DOCKER_OGEN_IMAGE} \
		-target ${SRC_DIR}/pkg/api/ogen/sample \
		-clean \
		${SRC_DIR}/api/sample/openapi_v3_ogen.yaml
.PHONY: openapi-server-ogen
