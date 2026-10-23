package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	// login route
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	// cambia nome utente
	rt.router.PUT("/users/me/name", rt.wrap(rt.setMyUserName))
	// search user route
	rt.router.GET("/users", rt.wrap(rt.searchUsers))
	// set Propic
	rt.router.PUT("/users/me/propic", rt.wrap(rt.setMyPhoto))
	// crea conversazione diretta
	rt.router.POST("/conversations", rt.wrap(rt.createConversation))
	// crea gruppo
	rt.router.POST("/groups", rt.wrap(rt.createGroup))
	// get conversations
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	// aggiungi al gruppo
	rt.router.POST("/groups/:id/members", rt.wrap(rt.addToGroup))
	// lista i membri del gruppo
	rt.router.GET("/groups/:id/members", rt.wrap(rt.getGroupMembers))
	// lascia il gruppo
	rt.router.DELETE("/groups/:id/members/me", rt.wrap(rt.leaveGroup))
	// imposta foto gruppo
	rt.router.PUT("/groups/:id/photo", rt.wrap(rt.setGroupPhoto))
	// imposta nome del gruppo
	rt.router.PUT("/groups/:id/name", rt.wrap(rt.setGroupName))
	// invia messaggio
	rt.router.POST("/conversations/:id/messages", rt.wrap(rt.sendMessage))
	// elimina messaggio
	rt.router.DELETE("/messages/:id", rt.wrap(rt.deleteMessage))
	// inoltra un messaggio
	rt.router.POST("/conversations/:id/messages/forward", rt.wrap(rt.forwardMessage))
	// aggiungi reaction
	rt.router.PUT("/messages/:id/reactions", rt.wrap(rt.addReaction))
	// rimuovi reazione
	rt.router.DELETE("/messages/:id/reactions", rt.wrap(rt.deleteReaction))
	// ricevi i messaggi di una conversazione
	rt.router.GET("/conversations/:id/messages", rt.wrap(rt.getConversation))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
