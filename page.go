package page_builder

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/slug"
	"github.com/qor/sorting"
	"github.com/qor/widget"
)

func init() {
	admin.RegisterViewPath("github.com/qor/page_builder/views")
}

type Page struct {
	gorm.Model
	Title         string
	TitleWithSlug slug.Slug

	QorWidgetSettings       []widget.QorWidgetSetting `gorm:"many2many:page_qor_widget_settings;ForeignKey:id;AssociationForeignKey:name"`
	QorWidgetSettingsSorter sorting.SortableCollection
}

func (*Page) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.UseTheme("page_builder")
	}
}
