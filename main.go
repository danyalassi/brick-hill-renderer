package main

import (
	"os"
	"sync"
	. "github.com/hawl1/brickgl"
	"github.com/nfnt/resize"
)

const (
	scale  = 4    // optional supersampling
	width  = 512 // output width in pixels
	height = 512 // output height in pixels
	fovy   = 74 // vertical field of view in degrees
	near   = 1    // near clipping plane
	far    = 100   // far clipping plane
)

var (
	//eye    = V(2.57161, 5.40692, 5.3927) // camera position
	//eye    = V(0, 0, 7)
	center = V(0, 3.13929, 0) // view center position
	//up     = V(0, 1, 0) // up vector
	eye    = V(1.5, 5, 4)
	//center = V(0, 3.5, 0)
	up     = V(0, 1, 0)
	//light  = V(-4.09796, -3.55779, 6.14444).Normalize() // light direction
	//light  = V(0.1, 0.45, 0.75).Normalize()
	light  = V(0.4, 0.7, 1.1).Normalize()
	//light  = V(0.75, 0.25, -0.5).Normalize()
	color  = HexColor("#468966") // object color
)

func main() {
	// load a mesh
	//if err != nil {
	//	panic(err)
	//}
	//mesh.BiUnitCube()
	//mesh.SmoothNormalsThreshold(Radians(30))
	//get vars
	headcolor := os.Args[1]
	torsocolor := os.Args[2]
	leftarmcolor := os.Args[3]
	rightarmcolor := os.Args[4]
	leftlegcolor := os.Args[5]
	rightlegcolor := os.Args[6]
	faceid := os.Args[7]
	headid := os.Args[8]
	tshirtid := os.Args[9]
	shirtid := os.Args[11]

	//allmeshes := [16]Mesh{}
	//allmeshes := [16]interface{}

	//create avatar

	
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	shader := NewPhongShader(matrix, light, eye)

	context := NewContext(width*scale, height*scale, 0, shader)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(){
		//context.ClearColorBufferWith(HexColor("#FFF8E3"))
		
		//LookAtDirection(forward, up Vector)
		mesh := NewEmptyObject()
		context.DrawObject(mesh, &wg)
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		headmesh := NewEmptyObject()
		//load head
		if headid == "0" {
			headmesh.AddMeshFromFile("Head.obj")
		}else{
			headmesh.AddMeshFromFile("/var/www/html/itemcache/"+headid+".obj")
		}
		//allmeshes[1] = mesh
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		shader = NewPhongShader(matrix, light, eye)
		headmesh.SetColor(HexColor(headcolor))
		//HexColor("#FEC400")
		//shader.SpecularPower = 0
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawObject(headmesh, &wg)
		wg.Done()

		facemesh := NewEmptyObject()
		//load face attempt
		if headid == "0" {
			facemesh.AddMeshFromFile("Head.obj")
		}else{
			facemesh.AddMeshFromFile("/var/www/html/itemcache/"+headid+".obj")
		}
		//allmeshes[2] = mesh
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		shader = NewPhongShader(matrix, light, eye)
		if faceid == "0" {
			facetexture, err := LoadTexture("face.png")
			if err != nil {
				panic(err)
			}
			headmesh.Texture = facetexture
		}else{
			facetexture, err := LoadTexture("/var/www/html/itemcache/"+faceid+".png")
			if err != nil {
				panic(err)
			}
			headmesh.Texture = facetexture
		}
		//HexColor("#FEC400")
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawObject(facemesh, &wg)
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		//load torso
		torsomesh := NewObjectFromFile("Torso.obj")
		//allmeshes[3] = mesh
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		shader = NewPhongShader(matrix, light, eye)
		torsomesh.SetColor(HexColor(torsocolor))
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawObject(torsomesh, &wg)
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		//load leftarm
		leftarmmesh := NewObjectFromFile("LeftArm.obj")
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		//allmeshes[4] = mesh
		shader = NewPhongShader(matrix, light, eye)
		leftarmmesh.SetColor(HexColor(leftarmcolor))
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawObject(leftarmmesh, &wg)
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		//load rightarm
		rightarmmesh := NewObjectFromFile("RightArm.obj")
		//allmeshes[5] = mesh
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		shader = NewPhongShader(matrix, light, eye)
		rightarmmesh.SetColor(HexColor(rightarmcolor))
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawObject(rightarmmesh, &wg)
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		if shirtid == "0" {	
			shirtid = "shirtStuds"
		}
		thisshatexture, err := LoadTexture("face.png")
		if err != nil {
			panic(err)
		}
		//load torso
		
		torsoshirtmesh := NewObjectFromFile("Torso.obj")
		//allmeshes[6] = mesh
		shader = NewPhongShader(matrix, light, eye)
		torsoshirtmesh.Texture = thisshatexture
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawObject(torsoshirtmesh, &wg)

		leftshirtmesh := NewObjectFromFile("LeftArm.obj")
		//allmeshes[7] = mesh
		shader = NewPhongShader(matrix, light, eye)
		leftshirtmesh.Texture = thisshatexture
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawObject(leftshirtmesh, &wg)

		rightshirtmesh := NewObjectFromFile("RightArm.obj")
		//allmeshes[8] = mesh
		shader = NewPhongShader(matrix, light, eye)
		rightshirtmesh.Texture = thisshatexture
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawObject(rightshirtmesh, &wg)
		wg.Done()
	}()
		/*mesh, _ = LoadOBJ("Shirt.obj")
		shader = NewPhongShader(matrix, light, eye)
		shader.Texture = thisshatexture
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(mesh)*/

		wg.Add(1)
		go func(){
			tshirtmesh := NewEmptyObject()
			//load tshirt if it exists
			if tshirtid != "0" {
				tshirtmesh.AddMeshFromFile("TShirt.obj")
				//allmeshes[9] = mesh
				shader = NewPhongShader(matrix, light, eye)
				tshtexture, err := LoadTexture("/var/www/html/itemcache/"+tshirtid+".png")
				if err != nil {
					panic(err)
				}
				tshirtmesh.Texture = tshtexture
				//shader.ObjectColor = HexColor(torsocolor)
				shader.DiffuseColor = Gray(0.25)
				shader.SpecularColor = Gray(0.1)
				shader.SpecularPower = 32
				shader.AmbientColor = HexColor("#c1bfbf")
				context.Shader = shader
				context.DrawObject(tshirtmesh, &wg)
			}
			wg.Done()
		}()

		wg.Add(1)
		go func(){
			//load leftleg
			leftlegmesh := NewObjectFromFile("LeftLeg.obj")
			//allmeshes[10] = mesh
			//mesh.SmoothNormals()
			//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
			shader = NewPhongShader(matrix, light, eye)
			leftlegmesh.SetColor(HexColor(leftlegcolor))
			shader.DiffuseColor = Gray(0.25)
			shader.SpecularColor = Gray(0.1)
			shader.SpecularPower = 32
			shader.AmbientColor = HexColor("#c1bfbf")
			context.Shader = shader
			context.DrawObject(leftlegmesh, &wg)
			wg.Done()
		}()
		
		wg.Add(1)
		go func(){
			//load rightleg
			rightlegmesh := NewObjectFromFile("RightLeg.obj")
			//allmeshes[11] = mesh
			//mesh.SmoothNormals()
			//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
			shader = NewPhongShader(matrix, light, eye)
			rightlegmesh.SetColor(HexColor(rightlegcolor))
			shader.DiffuseColor = Gray(0.25)
			shader.SpecularColor = Gray(0.1)
			shader.SpecularPower = 32
			shader.AmbientColor = HexColor("#c1bfbf")
			context.Shader = shader
			context.DrawObject(rightlegmesh, &wg)
			wg.Done()
		}()

	//if hatsplit != null{
        /*
		for _, element := range hatsplit {
			// element is the element from someSlice for where we are
			if element != "0" {
				mesh, _ = LoadOBJ("/var/www/html/itemcache/"+element+".obj")
				//allmeshes[len(allmeshes) + 1] = mesh
				//mesh.SmoothNormals()
				//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
				shader = NewPhongShader(matrix, light, eye)
				thishatexture, err := LoadTexture("/var/www/html/itemcache/"+element+".png")
				if err != nil {
					panic(err)
				}
				shader.Texture = thishatexture
				//shader.ObjectColor = HexColor(torsocolor)
				shader.DiffuseColor = Gray(0.25)
				shader.SpecularColor = Gray(0.1)
				shader.SpecularPower = 32
				shader.AmbientColor = HexColor("#c1bfbf")
				context.Shader = shader
				context.DrawMesh(mesh)
			}
		}
        */
	//}

	//matrix.LookAt(V(2, 5, 10), center, up)
	//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
	//.BoundingBox()


	// downsample image for antialiasing
	wg.Add(1)
	go func(){
		image := context.Image()
		image = resize.Resize(width, height, image, resize.Bilinear)

		// save image
		SavePNG("goout.png", image)
		wg.Done()
	}()
    }
