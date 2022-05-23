package testingpolicies_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/accuknox/auto-policy-discovery/src/types"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
	"log"
	"path/filepath"
	"strings"
)

func TestTestingpolicies(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testingpolicies Suite")
}

func ReadInstanceYaml(serverFile string, obj *types.KubeArmorPolicy)  {

	var files []string

	root := "/home/runner/work/discovery-engine/discovery-engine"
	//root := "tests"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		var res = strings.Contains(file, serverFile)
		if res == true {
			source, err1 := ioutil.ReadFile(file)

			 if err1 != nil {
			    log.Printf("Error: %v", err1.Error())
			 }

			 err1 = yaml.Unmarshal(source, &obj)
			 if err1 != nil {
			    log.Printf("Error: %v", err1.Error())
			 }
			break;
		}
	}
}

var _ = Describe("Knoxautopolicy validation", func(){
	Context("Checking Knoxautopolicy validation...", func(){
		It("return true", func(){
			f := types.KubeArmorPolicy{}

			ReadInstanceYaml("kubearmor_policies_default_explorer_knoxautopolicy", &f)

			Expect(f.Spec.Process.MatchPaths).NotTo(BeEmpty())
			Expect(f.Spec.File.MatchPaths).NotTo(BeEmpty())
			Expect(f.Spec.Network.MatchProtocols).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("Explorer Mysql validation", func(){
        Context("Checking Explorer Mysql validation...", func(){
                It("return true", func(){
                        f := types.KubeArmorPolicy{}

                        ReadInstanceYaml("kubearmor_policies_default_explorer_mysql", &f)

                        Expect(f.Spec.Process.MatchPaths).NotTo(BeEmpty())
                        Expect(f.Spec.File.MatchPaths).NotTo(BeEmpty())
                        Expect(f.Spec.Network.MatchProtocols).NotTo(BeEmpty())
                })
        })
})

var _ = Describe("Wordpress-Mysql Mysql validation", func(){
        Context("Checking Wordpress-Mysql Mysql validation...", func(){
                It("return true", func(){
                        f := types.KubeArmorPolicy{}

                        ReadInstanceYaml("kubearmor_policies_default_wordpress-mysql_mysql", &f)

                        Expect(f.Spec.Process.MatchPaths).NotTo(BeEmpty())
                        Expect(f.Spec.File.MatchPaths).NotTo(BeEmpty())
                        Expect(f.Spec.Network.MatchProtocols).NotTo(BeEmpty())
                })
        })
})

var _ = Describe("Wordpress-Mysql Wordpress validation", func(){
        Context("Checking Wordpress-Mysql Wordpress validation...", func(){
                It("return true", func(){
                        f := types.KubeArmorPolicy{}

                        ReadInstanceYaml("kubearmor_policies_default_wordpress-mysql_wordpress", &f)

                        Expect(f.Spec.Process.MatchPaths).NotTo(BeEmpty())
                        Expect(f.Spec.File.MatchPaths).NotTo(BeEmpty())
                        Expect(f.Spec.Network.MatchProtocols).NotTo(BeEmpty())
                })
        })
})
