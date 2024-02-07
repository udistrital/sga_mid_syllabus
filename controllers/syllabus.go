package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_mid_syllabus/utils"
	"github.com/udistrital/utils_oas/request"
)

// SyllabusController operations for Espacios_academicos
type SyllabusController struct {
	beego.Controller
}

// URLMapping ...
func (c *SyllabusController) URLMapping() {
	c.Mapping("PostSyllabusTemplate", c.PostSyllabusTemplate)
}

// PostSyllabusTemplate ...
// @Title PostSyllabusTemplate
// @Description post Syllabus template
// @Param   body        body    {}  true        "body generar plantilla del syllabus"
// @Success 200 {}
// @Failure 403 :body {}
// @router /generador-documentos [post]
func (c *SyllabusController) PostSyllabusTemplate() {
	var syllabusRequest map[string]interface{}
	var syllabusResponse map[string]interface{}
	var syllabusTemplateResponse map[string]interface{}
	var syllabusTemplateData map[string]interface{}
	var syllabusData map[string]interface{}

	failureAsn := map[string]interface{}{
		"Success": false,
		"Status":  "404",
		"Message": "Error service PostSyllabusTemplate: The request contains an incorrect parameter or no record exist",
		"Data":    nil}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &syllabusRequest); err == nil {
		syllabusCode := syllabusRequest["syllabusCode"]
		templateFormat, hasFormat := syllabusRequest["format"]
		if hasFormat {
			templateFormat = templateFormat.(string)
		} else {
			templateFormat = "pdf"
		}

		if syllabusVersion, hasVersion := syllabusRequest["version"]; hasVersion {
			syllabusErr := request.GetJson("http://"+beego.AppConfig.String("SyllabusService")+
				fmt.Sprintf("syllabus?query=syllabus_code:%v,version:%v&limit=1&offset=0", syllabusCode, syllabusVersion), &syllabusResponse)
			if syllabusErr != nil || syllabusResponse["Success"] == false {
				if syllabusErr == nil {
					syllabusErr = fmt.Errorf("SyllabusService: %v", syllabusResponse["Message"])
				}
				logs.Error(syllabusErr.Error())
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = syllabusErr.Error()
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			}
			syllabusList := syllabusResponse["Data"].([]interface{})
			if len(syllabusList) < 1 {
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = fmt.Errorf("SyllabusService: syllabus not found by syllabusCode and version")
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			} else {
				syllabusData = syllabusList[0].(map[string]interface{})
			}
		} else {
			syllabusErr := request.GetJson("http://"+beego.AppConfig.String("SyllabusService")+
				fmt.Sprintf("syllabus/%v", syllabusCode), &syllabusResponse)
			if syllabusErr != nil || syllabusResponse["Success"] == false {
				if syllabusErr == nil {
					syllabusErr = fmt.Errorf("SyllabusService: %v", syllabusResponse["Message"])
				}
				logs.Error(syllabusErr.Error())
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = syllabusErr.Error()
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			}
			syllabusData = syllabusResponse["Data"].(map[string]interface{})
		}

		spaceData, spaceErr := utils.GetAcademicSpaceData(
			int(syllabusData["plan_estudios_id"].(float64)),
			int(syllabusData["proyecto_curricular_id"].(float64)),
			int(syllabusData["espacio_academico_id"].(float64)))

		projectData, projectErr := utils.GetProyectoCurricular(int(syllabusData["proyecto_curricular_id"].(float64)))

		if spaceErr == nil && projectErr == nil {
			facultyData, facultyErr := utils.GetFacultadDelProyectoC(projectData["id_oikos"].(string))
			idiomas := ""

			if syllabusData["idioma_espacio_id"] != nil {
				idiomasStr, idiomaErr := utils.GetIdiomas(syllabusData["idioma_espacio_id"].([]interface{}))
				if idiomaErr == nil {
					idiomas = idiomasStr
				}
			}

			if facultyErr == nil {
				syllabusTemplateData = utils.GetSyllabusTemplateData(
					spaceData, syllabusData,
					facultyData, projectData, idiomas)

				utils.GetSyllabusTemplate(syllabusTemplateData, &syllabusTemplateResponse,
					fmt.Sprintf("%v", templateFormat))

				c.Data["json"] = map[string]interface{}{
					"Success": true,
					"Status":  "201",
					"Message": "Generated Syllabus Template OK",
					"Data":    syllabusTemplateResponse["Data"].(map[string]interface{})}
			} else {
				err := fmt.Errorf(
					"SyllabusTemplateService: Incomplete data to generate the document. Facultad y/o Idioma")
				logs.Error(err.Error())
				c.Ctx.Output.SetStatus(404)
				failureAsn["Data"] = err.Error()
				c.Data["json"] = failureAsn
				c.ServeJSON()
				return
			}
		} else {
			err := fmt.Errorf(
				"SyllabusTemplateService: Incomplete data to generate the document. Espacio Académico y/o Proyecto Curricular")
			logs.Error(err.Error())
			c.Ctx.Output.SetStatus(404)
			failureAsn["Data"] = err.Error()
			c.Data["json"] = failureAsn
			c.ServeJSON()
			return
		}
	}
	c.ServeJSON()
}
