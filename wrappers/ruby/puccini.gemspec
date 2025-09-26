version = ENV['PUCCINI_VERSION'] || '0.0.0'

Gem::Specification.new do |s|
  s.name                  = 'puccini'
  s.version               = version
  s.required_ruby_version = '>= 3.4'

  s.summary               = 'Puccini'
  s.description           = 'Deliberately stateless cloud topology management and deployment tools based on TOSCA'
  s.homepage              = 'https://github.com/tliron/go-puccini'
  s.license               = 'Apache-2.0'

  s.authors               = ['Tal Liron']
  s.email                 = 'tal.liron@gmail.com'

  s.files                 = [
  	'lib/puccini.rb',
  	'lib/libpuccini.so'
  ]
end
