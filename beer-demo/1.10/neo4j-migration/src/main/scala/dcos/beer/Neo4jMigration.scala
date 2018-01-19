package dcos.beer

import org.neo4j.driver.v1.{AuthTokens, GraphDatabase}
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.jdbc.datasource.DriverManagerDataSource

import scala.collection.JavaConversions._

@SpringBootApplication
object Neo4jMigration extends App {

  val mysqlDriver = new DriverManagerDataSource(sys.env("SQL_URL"))
  mysqlDriver.setDriverClassName("com.mysql.jdbc.Driver")
  val mysqlTemplate = new JdbcTemplate(mysqlDriver)

  val neo4jDriver = GraphDatabase.driver(sys.env("NEO4J_URL"), AuthTokens.basic(sys.env("NEO4J_USER"), sys.env("NEO4J_PASSWORD")))
  val neo4jSession = neo4jDriver.session()

  val beers = mysqlTemplate.queryForList("SELECT id, brewery_id, cat_id, style_id, name FROM `beers`").toList
  neo4jSession.run("CREATE " +
      beers.map(beer => {
        s"(: Beer {id: ${beer.get("id")}, " +
            s"brewery: ${beer.get("brewery_id")}, " +
            (if (beer.get("cat_id").toString != "-1") s"category: ${beer.get("cat_id")}, " else "") +
            (if (beer.get("style_id").toString != "-1") s"style: ${beer.get("style_id")}, " else "") +
            s"name: '${cleanString(beer.get("name"))}'})"
      }).mkString(",")
  )
  println("finished beers")

  val breweries = mysqlTemplate.queryForList("select b.id, b.name, b.city, b.state, b.country, g.latitude, g.longitude from `breweries` b LEFT JOIN `breweries_geocode` g ON b.id = g.brewery_id").toList
  neo4jSession.run("CREATE " +
      breweries.map(brewery => {
        s"(: Brewery {id: ${brewery.get("id")}, " +
            (if (brewery.get("state").toString != "") s"city: '${cleanString(brewery.get("city"))}', " else "") +
            (if (brewery.get("state").toString != "") s"state: '${cleanString(brewery.get("state"))}', " else "") +
            (if (brewery.get("state").toString != "") s"country: '${cleanString(brewery.get("country"))}', " else "") +
            s"latitude: ${brewery.get("latitude")}, " +
            s"longitude: ${brewery.get("longitude")}, " +
            s"name: '${cleanString(brewery.get("name"))}'})"
      }).mkString(",")
  )
  println("finished breweries")

  val categories = mysqlTemplate.queryForList("select id, cat_name from `categories`").toList
  neo4jSession.run("CREATE " +
      categories.map(category => {
        s"(: Category {id: ${category.get("id")}, " +
            s"name: '${cleanString(category.get("cat_name"))}'})"
      }).mkString(",")
  )
  println("finished categories")

  val styles = mysqlTemplate.queryForList("select id, cat_id, style_name from `styles`").toList
  neo4jSession.run("CREATE " +
      styles.map(style => {
        s"(: Style {id: ${style.get("id")}, " +
            s"category: ${style.get("cat_id")}, " +
            s"name: '${cleanString(style.get("style_name"))}'})"
      }).mkString(",")
  )
  println("finished styles")

  neo4jSession.run("CREATE INDEX ON :Beer(id)")
  neo4jSession.run("CREATE INDEX ON :Beer(name)")
  neo4jSession.run("CREATE INDEX ON :Brewery(id)")
  neo4jSession.run("CREATE INDEX ON :Brewery(name)")
  neo4jSession.run("CREATE INDEX ON :Category(id)")
  neo4jSession.run("CREATE INDEX ON :Style(id)")

  // create relations
  neo4jSession.run("MATCH (beer: Beer), (brewery: Brewery) WHERE beer.brewery = brewery.id MERGE (brewery)-[:PRODUCES]->(beer);")
  neo4jSession.run("MATCH (beer: Beer), (category: Category) WHERE beer.category = category.id MERGE (beer)-[:HAS_CATEGORY]->(category);")
  neo4jSession.run("MATCH (beer: Beer), (style: Style) WHERE beer.style = style.id MERGE (beer)-[:HAS_STYLE]->(style);")
  neo4jSession.run("MATCH (category: Category), (style: Style) WHERE style.category = category.id MERGE (style)-[:REFINES]->(category);")
  println("finished relations")
  println("finished neo4j migration")

  private def cleanString(s: AnyRef) = s.toString.replace("'", "`")
}
