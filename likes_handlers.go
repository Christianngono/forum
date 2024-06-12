package forum

import ( 
	"net/http"
	"strconv"
)

type Like struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

type Dislike struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

// Handler pour gérer les likes sur les posts
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err!= nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

	_, err = DB.Exec("UPDATE posts SET likes = likes + 1 WHERE id = ?", postID)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	http.Redirect(w, r, "/get-post?id="+ strconv.Itoa(postID), http.StatusSeeOther)
}


// Handler pour gérer les dislikes sur les posts
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
    if err!= nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    _, err = DB.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE id =?", postID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/get-post?id="+strconv.Itoa(postID), http.StatusSeeOther)	
}
