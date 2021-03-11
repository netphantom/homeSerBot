package web

import (
	"github.com/gin-contrib/multitemplate"
	"homeSerBot/pkg/forms"
	"homeSerBot/pkg/mysqlmodels"
)

type templateData struct {
	Form *forms.Form
	User *mysqlmodels.User
}

func templateRenderer() multitemplate.Renderer {
	mTemplate := multitemplate.New()
	mTemplate.AddFromFiles("home", "ui/template/baselayout.gohtml", "ui/template/home.gohtml", "ui/template/footer.gohtml")
	mTemplate.AddFromFiles("changePassword", "ui/template/baselayout.gohtml", "ui/template/password.gohtml", "ui/template/footer.gohtml")
	mTemplate.AddFromFiles("login", "ui/template/login.gohtml")
	mTemplate.AddFromFiles("profile", "ui/template/baselayout.gohtml", "ui/template/profile.gohtml", "ui/template/footer.gohtml")
	mTemplate.AddFromFiles("notifications", "ui/template/baselayout.gohtml", "ui/template/notifications.gohtml", "ui/template/footer.gohtml")
	mTemplate.AddFromFiles("process", "ui/template/baselayout.gohtml", "ui/template/processList.gohtml", "ui/template/footer.gohtml")
	mTemplate.AddFromFiles("processDetail", "ui/template/baselayout.gohtml", "ui/template/processMod.gohtml", "ui/template/footer.gohtml")
	mTemplate.AddFromFiles("processAdd", "ui/template/baselayout.gohtml", "ui/template/processAdd.gohtml", "ui/template/footer.gohtml")
	mTemplate.AddFromFiles("notFound", "ui/html/notfound.html")
	return mTemplate
}
