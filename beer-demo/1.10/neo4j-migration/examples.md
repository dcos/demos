# Query all styles

MATCH (brewery: Brewery)-[*1..2]->(style: Style) WHERE brewery.country = "Germany" WITH style, count(*) as c RETURN style, c;


# List all styles of beer
MATCH (brewery: Brewery)-[*1..2]->(style: Style) WHERE brewery.country = "Germany" WITH style, count(*) as c ORDER BY c DESC RETURN style.name, c;

# List all styles of beer produced in germany and count the breweries

MATCH (beer: Beer)<-[:PRODUCES]-(brewery: Brewery)-[*1..2]->(style: Style) WHERE brewery.country = "Germany" WITH style, count(beer) as c ORDER BY c DESC RETURN style.name, c;
