package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_syllabus_mid/controllers:SyllabusController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_syllabus_mid/controllers:SyllabusController"],
        beego.ControllerComments{
            Method: "PostSyllabusTemplate",
            Router: "/generador-documentos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_syllabus_mid/controllers:SyllabusLegacyController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_syllabus_mid/controllers:SyllabusLegacyController"],
        beego.ControllerComments{
            Method: "GetSyllabusLegacy",
            Router: "/syllabus/:qp_syllabus",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
