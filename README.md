# HiPeople Assignment - Image Store

This is an image store that allows the upload and retrieval of images using a simple API interface.

## Requirements

```go```

## How to run

Simply run ```go run main.go``` inside the root directory. API will be available at ```localhost:8080```

# API Endpoints

## Get image by name

Method: ```GET /?image=IMAGENAME.png```

### Examples

#### :star2: Found image by name :star2:

#### Request
```GET /?image=flying_bird.png```

#### Response

```Status 200 OK```

```Content-Type: application/octet-stream```

![flying bird](/flying_bird.png)


#### :boom: Image not found :boom:


#### Request

```GET /?image=ImageThatDoesNotExist.png```

#### Response

```Status 404 Not Found```

## Upload image
Method: ```POST /```

###### Maximum image size allowed is 32MB

### Examples

#### :star2: Successful image upload :star2:

##### Body: ```form-data```

```Key: image```

```Value: flying_bird.png (file)```

![flying bird](/flying_bird.png)

#### Response
```Status 200 OK```
```json
flying_bird.png
```

#### :boom: Invalid file extension :boom:

#### Request
```POST /```

##### Body: ```form-data```

```Key: image```

```Value: text_file.txt```
#### Response
```Status 400 Bad Request```
```json
file extension is not allowed
```

## License
[MIT](https://choosealicense.com/licenses/mit/)