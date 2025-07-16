# save to tmp file
cat openapi.json > openapi_tmp_in.json

# remove allOf from openapi_umlauts.json and save as openapi_no_allof.json
# jq '
# def walk(f):
#   . as $in
#   | if type == "object" then
#       reduce keys[] as $key
#         ( {}; . + { ($key): ($in[$key] | walk(f)) } ) | f
#     elif type == "array" then
#       map(walk(f)) | f
#     else
#       f
#     end;
#
# walk(
#   if type == "object" and has("allOf") and (.allOf | type == "array") and (.allOf[0]? | has("$ref")) then
#     del(.allOf) + { "$ref": .allOf[0]["$ref"] }
#   else .
#   end
# )
# ' openapi_tmp_in.json > openapi_tmp_out.json

# replace allOf with $ref
jq 'walk(
  if type == "object" and has("allOf") and (.allOf | type == "array") and (.allOf[0] | has("$ref"))
  then { "$ref": .allOf[0]["$ref"] }
  else .
  end
)' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# delete components.schemas.MailDomainRequest.properties.alias
jq 'del(.components.schemas.MailDomainRequest.properties.alias)' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# delete components.schemas.MailDomain.properties.alias
jq 'del(.components.schemas.MailDomain.properties.alias)' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# delete components.schemas.MailUserRequest.properties.alias
jq 'del(.components.schemas.MailUserRequest.properties.alias)' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# delete components.schemas.MailUser.properties.alias
jq 'del(.components.schemas.MailUser.properties.alias)' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# delete components.schemas.MailUserRequest.properties.domain
jq 'del(.components.schemas.MailUserRequest.properties.domain)' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# delete components.schemas.MailUser.properties.domain
jq 'del(.components.schemas.MailUser.properties.domain)' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# delete info.contact
jq 'del(.info.contact)' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# replace \u00fc\u00e4\u00f6\u00dc\u00c4\u00d6\u00df with üäöÜÄÖß in openapi.json and save as openapi_fixed.json
sed 's/\\\\u00fc/ü/g; s/\\\\u00e4/ä/g; s/\\\\u00f6/ö/g; s/\\\\u00dc/Ü/g; s/\\\\u00c4/Ä/g; s/\\\\u00d6/Ö/g; s/\\\\u00df/ß/g' openapi_tmp_in.json > openapi_tmp_out.json
cat openapi_tmp_out.json > openapi_tmp_in.json

# save the final cleaned up openapi.json
cat openapi_tmp_in.json > openapi_clean.json

# generate provider code and client code
echo "Generating provider code and client code..."
tfplugingen-openapi generate \
  --config generator_config.yml \
  --output provider_code_spec.json \
  openapi_clean.json

# generate provider and client code from the provider_code_spec.json
echo "Generating provider and client code from provider_code_spec.json..."
tfplugingen-framework generate all \
  --input provider_code_spec.json \
  --output gen/provider

# generate client code
# go tool ogen --package client --target gen/client --clean openapi_clean.json