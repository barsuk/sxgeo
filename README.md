# sxgeo
A tool to work with Sypex Geo database, which helps detect a country and a city by IP

Надстройка над базой данных Sypex Geo версии 2.2 (https://sypexgeo.net/ru/docs/)

## Настройка
В базу записаны коды, зависящие от машинного порядка записи байтов (LittleEndian, BigEndian).
Для начала работы нужно определить порядок на рабочей машине.
По умолчанию в переменную hbo установлена LittleEndian. 
```
var hbo = binary.LittleEndian
```
Кодировку определяйте с помощью функции DetectEndian(), задавайте SetEndian(sxgeo.LITTLE || sxgeo.BIG).


