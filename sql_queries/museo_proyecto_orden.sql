SELECT * FROM museo_proyecto.orden;

CREATE EVENT IF NOT EXISTS cancelar_despues_de_2h
ON SCHEDULE EVERY 1 HOUR
DO
  UPDATE orden
  SET status = 'CANCELADA'
  WHERE status = 'PENDIENTE' AND STR_TO_DATE(fecha, '%Y-%m-%dT%H:%i:%s.%fZ') < (NOW() - INTERVAL 2 HOUR);

-- Specifier	Meaning	Example
-- %Y	4-digit year	2026
-- %y	2-digit year	26
-- %m	Month (01–12)	03
-- %d	Day (01–31)	07
-- %H	Hour (00–23)	14
-- %i	Minutes (00–59)	45
-- %s	Seconds (00–59)	30