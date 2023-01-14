//файл cmd/client/client.go
package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pablogolobaro/pdfcomposer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"path/filepath"

	"time"
)

func main() {

	cwt, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(cwt, "localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	pdfComposeClient := pdfcomposer.NewPdfComposeClient(conn)
	imagePaths := []string{"./photo1.jpg", "./photo2.jpg"}
	uploadImage(pdfComposeClient, imagePaths)
}

func uploadImage(client pdfcomposer.PdfComposeClient, imagePaths []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := client.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}
	req := &pdfcomposer.UploadImageRequest{
		Data: &pdfcomposer.UploadImageRequest_Count{
			Count: int32(len(imagePaths)),
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}
	for _, imagePath := range imagePaths {
		err = sendImage(stream, imagePath)
		if err != nil {
			log.Fatal("cannot upload image: ", err)
		}
	}
	err = stream.CloseSend()
	if err != nil {
		log.Fatal("cannot close stream: ", err)
	}
	err = getPDF(stream)
	if err != nil {
		log.Fatal(err)
	}
}
func sendImage(stream pdfcomposer.PdfCompose_UploadImageClient, imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()
	req := &pdfcomposer.UploadImageRequest{
		Data: &pdfcomposer.UploadImageRequest_Info{
			Info: &pdfcomposer.ImageInfo{
				ImageType: filepath.Ext(imagePath),
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pdfcomposer.UploadImageRequest{
			Data: &pdfcomposer.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}
	req = &pdfcomposer.UploadImageRequest{
		Data: &pdfcomposer.UploadImageRequest_EndFile{
			EndFile: true,
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}
	res, err := stream.Recv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}
	log.Printf("Success: %v, image uploaded with size: %d", res.GetInfo().GetSuccess(), res.GetInfo().GetSize())
	return nil
}
func getPDF(stream pdfcomposer.PdfCompose_UploadImageClient) error {
	log.Print("Getting pdf file...")
	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		log.Print("waiting to receive more data")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}
		chunk := req.GetChunkData()
		size := len(chunk)
		imageSize += size
		log.Printf("received a chunk with size: %d", size)
		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}
	log.Print("Got pdf file...Saving...")
	_, err := savePDF(imageData)
	if err != nil {
		return err
	}
	return nil
}
func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
func savePDF(imageData bytes.Buffer) (io.ReadCloser, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", ".", imageID, ".pdf")

	file, err := os.Create(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create image file: %w", err)
	}
	defer file.Close()
	_, err = imageData.WriteTo(file)
	if err != nil {
		return nil, fmt.Errorf("cannot write pdf to file: %w", err)
	}
	return file, nil
}
