package page_builder

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/slug"
	"github.com/qor/widget"
)

func init() {
	admin.RegisterViewPath("github.com/qor/page_builder/views")
}

type Page struct {
	gorm.Model
	Title         string
	TitleWithSlug slug.Slug

	Containers Containers `sql:"type:text;"`
}

func (*Page) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		res.UseTheme("page_builder")
	}
}

type Container struct {
	Name string
}

type Containers []Container

func (containers *Containers) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, containers)
	case string:
		if v != "" {
			return containers.Scan([]byte(v))
		}
	case []string:
		c := Containers{}
		for _, name := range v {
			if name != "" {
				c = append(c, Container{Name: name})
			}
		}
		*containers = c
	default:
		return errors.New("not supported")
	}

	return nil
}

func (containers Containers) Value() (driver.Value, error) {
	if len(containers) == 0 {
		return nil, nil
	}

	return json.Marshal(containers)
}

func (page Page) ContainerRecords(db *gorm.DB) (records []widget.QorWidgetSetting) {
	names := []string{}
	for _, container := range page.Containers {
		names = append(names, container.Name)
	}

	containers := []widget.QorWidgetSetting{}
	db.Where("name in (?)", names).Find(&containers)

	for _, name := range names {
		for _, container := range containers {
			if container.Name == name {
				records = append(records, container)
			}
		}
	}

	return
}
