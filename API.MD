collection:
	name
	data: []map[NamedType][]*type from NamedType*
NamedType:
	name: unique
	type: of data

Basic Types:
	num:
		timestamp
		value: float or smt...
	str:
		timestamp
		value: string


api:
	collection/:
		GET: collections
		POST: new collection
	collection/{id}:
		GET: collection
		POST: add data
		DELETE: collection
	collection/data/{id_c}/{id_d}:
		GET
		DELETE
	type/named:
		GET
		POST: new type
	type/{name}:
		GET
		DELETE
	type/basic:
		GET
