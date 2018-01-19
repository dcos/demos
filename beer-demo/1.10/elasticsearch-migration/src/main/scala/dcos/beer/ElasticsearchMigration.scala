package dcos.beer

import java.net.InetAddress

import org.elasticsearch.client.Client
import org.elasticsearch.common.settings.Settings
import org.elasticsearch.common.transport.InetSocketTransportAddress
import org.elasticsearch.common.xcontent.XContentFactory.jsonBuilder
import org.elasticsearch.transport.client.PreBuiltTransportClient
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.jdbc.datasource.DriverManagerDataSource

import scala.collection.JavaConversions._
import scala.collection.JavaConverters._

@SpringBootApplication
object ElasticsearchMigration extends App {

  val index = "beer"

  val mysqlDriver = new DriverManagerDataSource(sys.env("SQL_URL"))
  mysqlDriver.setDriverClassName("com.mysql.jdbc.Driver")
  val mysqlTemplate = new JdbcTemplate(mysqlDriver)

  val settings = Settings.builder() /*.put("cluster.name", "myClusterName")*/ .build()
  val client: Client = new PreBuiltTransportClient(settings)
      .addTransportAddress(new InetSocketTransportAddress(InetAddress.getByName(sys.env("ELASTICSEARCH_URL")), 9300))

  if (!client.admin().indices().prepareExists(index).execute().actionGet().isExists) {
    val createResponse = client.admin().indices()
        .prepareCreate(index)
        .addMapping("beer", jsonBuilder()
            .startObject()
            .field("properties", Map(
              "name" -> Map("type" -> "text", "store" -> "yes").asJava,
              "description" -> Map("type" -> "text", "store" -> "yes").asJava
            ).asJava)
            .endObject()
        ).addMapping("brewery", jsonBuilder()
        .startObject()
        .field("properties", Map(
          "name" -> Map("type" -> "text", "store" -> "yes").asJava,
          "description" -> Map("type" -> "text", "store" -> "yes").asJava
        ).asJava)
        .endObject()
    )
        .execute().actionGet()
  }
  val beers = mysqlTemplate.queryForList("SELECT id, name, descript FROM `beers` where descript != ''").toList
  val bulkIndex = client.prepareBulk()
  beers.map {
    beer =>
      client.prepareIndex(index, "beer")
          .setSource(
            jsonBuilder()
                .startObject()
                .field("name", s"${beer.get("name")}")
                .field("description", s"${beer.get("descript")}")
                .endObject()
          ).setCreate(true).setId(s"${beer.get("id")}")
  }.foreach(bulkIndex.add)
  println("finished beers")

  val breweries = mysqlTemplate.queryForList("select id, name, descript from `breweries` where descript != ''").toList
  breweries.map {
    brewery =>
      client.prepareIndex(index, "brewery")
          .setSource(
            jsonBuilder()
                .startObject()
                .field("name", s"${brewery.get("name")}")
                .field("description", s"${brewery.get("descript")}")
                .endObject()
          ).setCreate(true).setId(s"${brewery.get("id")}")
  }.foreach(bulkIndex.add)
  println("finished breweries")

  val response = bulkIndex.execute().actionGet()
  println(s"finished elasticsearch migration: ${response.buildFailureMessage()}")



  private def cleanString(s: AnyRef) = s.toString.replace("'", "`")
}
