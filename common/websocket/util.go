package websocket

func GenGroupKey(systemId, groupName string) string {
	return systemId + ":" + groupName
}
