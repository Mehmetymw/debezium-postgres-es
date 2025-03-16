docker-compose down && docker-compose up -d --build

sleep 10

# Orders tablosunu oluştur
docker exec -i debezium-postgres-1 psql -U postgres inventory -c "CREATE TABLE IF NOT EXISTS public.orders (
    id VARCHAR(255) PRIMARY KEY, 
    orderId VARCHAR(255), 
    customerId VARCHAR(255), 
    status VARCHAR(50),
    location GEOMETRY(Point, 4326)
);"

# Mevcut verileri ekle (geometri kolonu ile birlikte)
docker exec -i debezium-postgres-1 psql -U postgres inventory -c "
INSERT INTO public.orders (id, orderId, customerId, status, location) VALUES 
('1', '101', '501', 'NEW', ST_SetSRID(ST_MakePoint(28.9784, 41.0082), 4326)),
('2', '102', '502', 'PROCESSING', ST_SetSRID(ST_MakePoint(28.9795, 41.0256), 4326)),
('3', '103', '503', 'COMPLETED', ST_SetSRID(ST_MakePoint(29.0469, 41.0128), 4326)),
('4', '104', '504', 'SHIPPED', ST_SetSRID(ST_MakePoint(28.9651, 41.0162), 4326)),
('5', '105', '505', 'DELIVERED', ST_SetSRID(ST_MakePoint(28.9833, 41.0351), 4326)),
('6', '106', '506', 'CANCELLED', ST_SetSRID(ST_MakePoint(29.0225, 41.0183), 4326)),
('7', '107', '507', 'RETURNED', ST_SetSRID(ST_MakePoint(28.9917, 41.0422), 4326)),
('8', '108', '508', 'PENDING', ST_SetSRID(ST_MakePoint(29.0107, 41.0890), 4326)),
('9', '109', '509', 'ON_HOLD', ST_SetSRID(ST_MakePoint(29.0357, 41.0466), 4326)),
('10', '110', '510', 'BACKORDERED', ST_SetSRID(ST_MakePoint(28.9530, 41.0234), 4326));"

# Connectorları sil (eğer varsa)
curl -X DELETE http://localhost:8083/connectors/postgres-source || echo "Connector not found"
curl -X DELETE http://localhost:8083/connectors/elastic-sink || echo "Connector not found"

# Connectorları yapılandır
curl -X POST -H "Content-Type: application/json" --data @postgres-source.json http://localhost:8083/connectors
curl -X POST -H "Content-Type: application/json" --data @elastic-sink.json http://localhost:8083/connectors
curl -X POST -H "Content-Type: application/json" --data @postgres-source-all.json http://localhost:8083/connectors
curl -X POST -H "Content-Type: application/json" --data @elastic-sink-all.json http://localhost:8083/connectors


# Elasticsearch'te verileri kontrol et
echo "İlk senkronizasyon sonrası Elasticsearch'teki veriler:"
curl -X GET "http://localhost:9200/dbserver1.public.orders/_search?pretty"

# Connector durumlarını kontrol et
echo "PostgreSQL connector durumu:"
curl -X GET http://localhost:8083/connectors/postgres-source/status

echo "Elasticsearch connector durumu:"
curl -X GET http://localhost:8083/connectors/elastic-sink/status
