package cupcake

import (
	"net/http"
	"reflect"

	"github.com/lz-nsc/cupcake/log"
	"github.com/lz-nsc/cupcake/orm/session"
)

type Controller interface {
	Create(*Response, *Request)
	Retrive(*Response, *Request)
	Update(*Response, *Request)
	Delete(*Response, *Request)
}

type BaseController struct {
	Model   interface{}
	session *session.Session
	// TODO: might need to move to serializer
	fields []string
}

func NewBaseController(model interface{}) *BaseController {
	session := newDBSession()
	err := session.Model(model)
	if err != nil {
		log.Errorf("Failed to create db session with model %T, err: %s", model, err)
		return nil
	}

	fields := []string{}
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()

	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)
		fields = append(fields, f.Name)
	}
	return &BaseController{
		Model:   model,
		session: session,
		fields:  fields,
	}
}

func (base *BaseController) Create(resp *Response, req *Request) {
	if base.Model == nil {
		log.Error("model cannot be nil in controller")
		resp.Error(http.StatusInternalServerError, "Internal Error")
		return
	}

	// Check whether table exist in database
	if exists := base.session.HasTable(); !exists {
		// If not, create table
		err := base.session.CreateTable()
		if err != nil {
			log.Errorf("failed to create table for model %s", base.session.ModelName())
			resp.Error(http.StatusInternalServerError, "Internal Error")
		}
	}
	instance := reflect.New(reflect.Indirect(reflect.ValueOf(base.Model)).Type()).Interface()

	// Parse request data
	// TODO: Currently can only accept JSON, whether I should consider how to parse Postform
	err := req.Parse(instance)

	if err != nil {
		log.Errorf("failed to read request body, err: %s\n", err.Error())
		resp.Error(http.StatusBadRequest, "Bad Request")
		return
	}

	// Insert new data to database
	count, err := base.session.Insert(instance)
	if err != nil {
		log.Errorf("failed to insert record, err: %s\n", err.Error())
		resp.Error(http.StatusInternalServerError, "Internal Error")
		return
	}
	log.Infof("Successfully insert %d row(s)\n", count)
	resp.Status(http.StatusCreated)
}

// Retrive data with givn primary key
func (base *BaseController) Retrive(resp *Response, req *Request) {
	pk := req.Param("pk")
	instance := reflect.New(reflect.Indirect(reflect.ValueOf(base.Model)).Type()).Interface()
	err := base.session.FindOneWithPK(pk, instance)
	if err != nil {
		resp.Error(http.StatusBadRequest, err.Error())
		return
	}
	resp.JSON(http.StatusOK, instance)
}
func (base *BaseController) Update(resp *Response, req *Request) {
	resp.Error(http.StatusMethodNotAllowed, "Method Not Allowed")
}
func (base *BaseController) Delete(resp *Response, req *Request) {
	resp.Error(http.StatusMethodNotAllowed, "Method Not Allowed")
}

func (base *BaseController) Parse(req *Request, obj interface{}) {
	req.Parse(obj)
}
