package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/pablogolobaro/pdfcomposer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"log"
	"os"
	"service-pdf-compose-grpc/pkg/composer"
)

type Server struct {
	pb.UnimplementedPdfComposeServer
}

func (server *Server) UploadImage(stream pb.PdfCompose_UploadImageServer) error {
	pdf, err := getImagesAndConverToPDF(stream)
	if err != nil {
		return err
	}
	err = sendPDF(stream, pdf)
	if err != nil {
		return err
	}
	return nil
}
func getImage(stream pb.PdfCompose_UploadImageServer) (io.ReadCloser, error) {
	req, err := stream.Recv()
	if err != nil {
		return nil, logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload-image request with image type %s", imageType)
	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		log.Print("waiting to receive more data")

		req, err := stream.Recv()
		if v := req.GetEndFile(); v == true {
			log.Print("no more data")
			break
		}
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return nil, logError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}
		chunk := req.GetChunkData()
		size := len(chunk)
		imageSize += size
		log.Printf("received a chunk with size: %d", size)
		_, err = imageData.Write(chunk)
		if err != nil {
			return nil, logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}

	}
	res := &pb.UploadImageResponse{
		Data: &pb.UploadImageResponse_Info{Info: &pb.UploadResult{
			Success: true,
			Size:    uint32(imageSize),
		}},
	}
	err = stream.Send(res)
	if err != nil {
		return nil, logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}
	log.Printf("got image with size: %d", imageSize)
	readCloser := ioutil.NopCloser(&imageData)
	return readCloser, nil
}

func getImagesAndConverToPDF(stream pb.PdfCompose_UploadImageServer) (io.ReadCloser, error) {
	files := []io.ReadCloser{}
	req, err := stream.Recv()
	if err != nil {
		return nil, logError(status.Errorf(codes.Unknown, "cannot receive image count"))
	}
	imageCount := req.GetCount()
	for i := 0; i < int(imageCount); i++ {
		file, err := getImage(stream)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	pdfFile, err := composer.ComposeFromFiles(files)
	if err != nil {
		return nil, logError(status.Errorf(codes.Unknown, "cannot convert image: %v", err))
	}
	log.Printf("File has been converted ")
	return pdfFile, nil
}

func sendPDF(stream pb.PdfCompose_UploadImageServer, pdf io.ReadCloser) error {
	defer pdf.Close()
	reader := bufio.NewReader(pdf)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadImageResponse{
			Data: &pb.UploadImageResponse_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}
	log.Print("Sent PDF File")
	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

func saveImage(imageType string, imageData bytes.Buffer) (io.ReadCloser, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", ".", imageID, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create image file: %w", err)
	}
	defer file.Close()
	_, err = imageData.WriteTo(file)
	if err != nil {
		return nil, fmt.Errorf("cannot write image to file: %w", err)
	}
	return file, nil
}
