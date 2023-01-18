Pod::Spec.new do |spec|
  spec.name         = 'nbn'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/token'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS nbn Client'
  spec.source       = { :git => 'https://github.com/token.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/nbn.framework'

	spec.prepare_command = <<-CMD
    curl https://nbnstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/nbn.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
