package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["sga_mid_syllabus/controllers:SyllabusController"] = append(beego.GlobalControllerRouter["sga_mid_syllabus/controllers:SyllabusController"],
        beego.ControllerComments{
            Method: "PostSyllabusTemplate",
            Router: "/generador-documentos",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sga_mid_syllabus/controllers:SyllabusLegacyController"] = append(beego.GlobalControllerRouter["sga_mid_syllabus/controllers:SyllabusLegacyController"],
        beego.ControllerComments{
            Method: "GetSyllabusLegacy",
            Router: "/syllabus/:qp_syllabus",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
