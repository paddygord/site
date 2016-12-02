function toggle_autoupdate() {
}

function set_freshness(value) {
    var d = new Date();
    d.setTime(d.getTime() + (7*24*60*60*1000));
    document.cookie = "freshness=" + value + ";expires=" + d.toUTCString() + ";path=/;Secure";
}

function toggle_hidden(element) {
    if (element.innerHTML == "hide") {
        element.innerHTML = "show"
        els = element.parentNode.querySelectorAll("*")
        for (var i = 0; i < els.length; i++) {
            if (els[i] != element) {
                els[i].style.display = "none"
            }
        }
    } else {
        els = element.parentNode.querySelectorAll("*")
        for (var i = 0;i < els.length; i++) {
            if (els[i] != element) {
                els[i].style.display = null
            }
        }
        element.innerHTML = "hide"
    }
}

function toggle_comment(element) {
    var comment = document.querySelector("#comment")
    if (comment == element.nextSibling) {
        comment.style.display = null
        e = document.querySelector("#post")
        e.parentElement.insertBefore(comment, e.nextSibling)
    } else {
        comment.style.display = "initial"
        post = document.querySelector(".post")
        comment.post_id.value = post.id
        if (element.parentElement != post) {
            comment.comment_id.value = element.parentElement.id
        }
        element.parentElement.insertBefore(comment, element.nextSibling)
    }
}