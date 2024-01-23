package main

import (
	"os"
	. "github.com/fogleman/fauxgl"
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

	mesh := NewEmptyMesh()

	context := NewContext(width*scale, height*scale)
	//context.ClearColorBufferWith(HexColor("#FFF8E3"))
	
	//LookAtDirection(forward, up Vector)
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	shader := NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader
	context.DrawMesh(mesh)

	headmesh := NewEmptyMesh()
	//load head
	if headid == "0" {
		headmesh, _ = LoadOBJ("Head.obj")
	}else{
		headmesh, _ = LoadOBJ("/var/www/html/itemcache/"+headid+".obj")
	}
	//allmeshes[1] = mesh
	//mesh.SmoothNormals()
	//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
	shader = NewPhongShader(matrix, light, eye)
	shader.ObjectColor = HexColor(headcolor)
	//HexColor("#FEC400")
	//shader.SpecularPower = 0
	shader.DiffuseColor = Gray(0.25)
	shader.SpecularColor = Gray(0.1)
	shader.SpecularPower = 32
	shader.AmbientColor = HexColor("#c1bfbf")
	context.Shader = shader
	context.DrawMesh(headmesh)

	facemesh := NewEmptyMesh()
	//load face attempt
	if headid == "0" {
		facemesh, _ = LoadOBJ("Head.obj")
	}else{
		facemesh, _ = LoadOBJ("/var/www/html/itemcache/"+headid+".obj")
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
		shader.Texture = facetexture
	}else{
		facetexture, err := LoadTexture("/var/www/html/itemcache/"+faceid+".png")
		if err != nil {
			panic(err)
		}
		shader.Texture = facetexture
	}
	//HexColor("#FEC400")
	shader.DiffuseColor = Gray(0.25)
	shader.SpecularColor = Gray(0.1)
	shader.SpecularPower = 32
	shader.AmbientColor = HexColor("#c1bfbf")
	context.Shader = shader
	context.DrawMesh(facemesh)

	//load torso
	torsomesh, _ := LoadOBJ("Torso.obj")
	//allmeshes[3] = mesh
	//mesh.SmoothNormals()
	//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
	shader = NewPhongShader(matrix, light, eye)
	shader.ObjectColor = HexColor(torsocolor)
	shader.DiffuseColor = Gray(0.25)
	shader.SpecularColor = Gray(0.1)
	shader.SpecularPower = 32
	shader.AmbientColor = HexColor("#c1bfbf")
	context.Shader = shader
	context.DrawMesh(torsomesh)

		//load leftarm
		leftarmmesh, _ := LoadOBJ("LeftArm.obj")
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		//allmeshes[4] = mesh
		shader = NewPhongShader(matrix, light, eye)
		shader.ObjectColor = HexColor(leftarmcolor)
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(leftarmmesh)
	
		//load rightarm
		rightarmmesh, _ := LoadOBJ("RightArm.obj")
		//allmeshes[5] = mesh
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		shader = NewPhongShader(matrix, light, eye)
		shader.ObjectColor = HexColor(rightarmcolor)
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(rightarmmesh)
	
		if shirtid == "0" {	
			shirtid = "shirtStuds"
		}
		thisshatexture, err := LoadTexture("face"+".png")
		if err != nil {
			panic(err)
		}
		//load torso
		
		torsoshirtmesh, _ := LoadOBJ("Torso.obj")
		//allmeshes[6] = mesh
		shader = NewPhongShader(matrix, light, eye)
		shader.Texture = thisshatexture
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(torsoshirtmesh)
	
		leftshirtmesh, _ := LoadOBJ("LeftArm.obj")
		//allmeshes[7] = mesh
		shader = NewPhongShader(matrix, light, eye)
		shader.Texture = thisshatexture
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(leftshirtmesh)
	
		rightshirtmesh, _ := LoadOBJ("RightArm.obj")
		//allmeshes[8] = mesh
		shader = NewPhongShader(matrix, light, eye)
		shader.Texture = thisshatexture
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(rightshirtmesh)
	
		/*mesh, _ = LoadOBJ("Shirt.obj")
		shader = NewPhongShader(matrix, light, eye)
		shader.Texture = thisshatexture
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(mesh)*/
		
		tshirtmesh := NewEmptyMesh()
		//load tshirt if it exists
		if tshirtid != "0" {
			tshirtmesh, _ = LoadOBJ("TShirt.obj")
			//allmeshes[9] = mesh
			shader = NewPhongShader(matrix, light, eye)
			tshtexture, err := LoadTexture("/var/www/html/itemcache/"+tshirtid+".png")
			if err != nil {
				panic(err)
			}
			shader.Texture = tshtexture
			//shader.ObjectColor = HexColor(torsocolor)
			shader.DiffuseColor = Gray(0.25)
			shader.SpecularColor = Gray(0.1)
			shader.SpecularPower = 32
			shader.AmbientColor = HexColor("#c1bfbf")
			context.Shader = shader
			context.DrawMesh(tshirtmesh)
		}
			
		//load leftleg
		leftlegmesh, _ := LoadOBJ("LeftLeg.obj")
		//allmeshes[10] = mesh
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		shader = NewPhongShader(matrix, light, eye)
		shader.ObjectColor = HexColor(leftlegcolor)
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(leftlegmesh)
				
		//load rightleg
		rightlegmesh, _ := LoadOBJ("RightLeg.obj")
		//allmeshes[11] = mesh
		//mesh.SmoothNormals()
		//mesh.Transform(Scale(V(2.5, 2.5, 2.5)))
		shader = NewPhongShader(matrix, light, eye)
		shader.ObjectColor = HexColor(rightlegcolor)
		shader.DiffuseColor = Gray(0.25)
		shader.SpecularColor = Gray(0.1)
		shader.SpecularPower = 32
		shader.AmbientColor = HexColor("#c1bfbf")
		context.Shader = shader
		context.DrawMesh(rightlegmesh)

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
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)

	// save image
	SavePNG("goout.png", image)
    }
