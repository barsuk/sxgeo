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

## Использование
Cчитайте файл SxGeoCity.dat в память:
```
	if _, err := sxgeo.ReadDBToMemory(dbPath); err != nil {
		log.Fatalf("error: cannot read database file: %v", err)
	}
```
IP строкой вида 8.8.8.8 передайте в функцию GetCityFull:
```
		city, err := sxgeo.GetCityFull(ip)
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}
```
Теперь можно преобразовать полученную структуру, например, в json. 
```
		enc, err := json.Marshal(city)
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", enc)
```
