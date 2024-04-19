package services

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/sga_syllabus_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func PostSyllabusTemplate(data []byte) requestresponse.APIResponse {
	var syllabusRequest map[string]interface{}
	var syllabusResponse map[string]interface{}
	var syllabusTemplateResponse map[string]interface{}
	var syllabusTemplateData map[string]interface{}
	var syllabusData map[string]interface{}

	if err := json.Unmarshal(data, &syllabusRequest); err == nil {
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
				return requestresponse.APIResponseDTO(false, 404, nil, syllabusErr.Error())
			}
			syllabusList := syllabusResponse["Data"].([]interface{})
			if len(syllabusList) < 1 {
				return requestresponse.APIResponseDTO(false, 404, nil, fmt.Errorf("SyllabusService: syllabus not found by syllabusCode and version"))
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
				return requestresponse.APIResponseDTO(false, 404, nil, syllabusErr.Error())
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

				return requestresponse.APIResponseDTO(true,
					201, syllabusTemplateResponse["body"].(map[string]interface{}),
					"Generated Syllabus Template OK")
			} else {
				err := fmt.Errorf(
					"SyllabusTemplateService: Incomplete data to generate the document. Facultad y/o Idioma")
				logs.Error(err.Error())
				return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
			}
		} else {
			err := fmt.Errorf(
				"SyllabusTemplateService: Incomplete data to generate the document. Espacio Académico y/o Proyecto Curricular")
			logs.Error(err.Error())
			return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
		}
	}
	return requestresponse.APIResponseDTO(false, 400, nil, "Invalid request body")
}
