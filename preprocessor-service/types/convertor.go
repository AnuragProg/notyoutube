package types

import (

	"github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/database/postgres"
	dbType "github.com/anuragprog/notyoutube/preprocessor-service/types/database"
	mqType "github.com/anuragprog/notyoutube/preprocessor-service/types/mq"
)

var WorkerTypeToPostgresWorkerType = map[dbType.WorkerType]postgres.WorkerType{
	dbType.WorkerTypeVideoEncoder:       postgres.WorkerTypeVideoEncoder,
	dbType.WorkerTypeAsciiEncoder:       postgres.WorkerTypeAsciiEncoder,
	dbType.WorkerTypeThumbnailGenerator: postgres.WorkerTypeThumbnailGenerator,
	dbType.WorkerTypeAssembler:          postgres.WorkerTypeAssembler,
	dbType.WorkerTypeVideoExtractor:     postgres.WorkerTypeVideoExtractor,
	dbType.WorkerTypeAudioExtractor:     postgres.WorkerTypeAudioExtractor,
	dbType.WorkerTypeMetadataExtractor:  postgres.WorkerTypeMetadataExtractor,
}
var PostgresWorkerTypeToWorkerType = map[postgres.WorkerType]dbType.WorkerType{
	postgres.WorkerTypeVideoEncoder:       dbType.WorkerTypeVideoEncoder,
	postgres.WorkerTypeAsciiEncoder:       dbType.WorkerTypeAsciiEncoder,
	postgres.WorkerTypeThumbnailGenerator: dbType.WorkerTypeThumbnailGenerator,
	postgres.WorkerTypeAssembler:          dbType.WorkerTypeAssembler,
	postgres.WorkerTypeVideoExtractor:     dbType.WorkerTypeVideoExtractor,
	postgres.WorkerTypeAudioExtractor:     dbType.WorkerTypeAudioExtractor,
	postgres.WorkerTypeMetadataExtractor:  dbType.WorkerTypeMetadataExtractor,
}
var ProtoWorkerTypeToPostgresWorkerType = map[mqType.WorkerType]postgres.WorkerType{
	mqType.WorkerType_VIDEO_ENCODER:       postgres.WorkerTypeVideoEncoder,
	mqType.WorkerType_ASCII_ENCODER:       postgres.WorkerTypeAsciiEncoder,
	mqType.WorkerType_THUMBNAIL_GENERATOR: postgres.WorkerTypeThumbnailGenerator,
	mqType.WorkerType_ASSEMBLER:           postgres.WorkerTypeAssembler,
	mqType.WorkerType_VIDEO_EXTRACTOR:     postgres.WorkerTypeVideoExtractor,
	mqType.WorkerType_AUDIO_EXTRACTOR:     postgres.WorkerTypeAudioExtractor,
	mqType.WorkerType_METADATA_EXTRACTOR:  postgres.WorkerTypeMetadataExtractor,
}
var ProtoWorkerTypeToWorkerType = map[mqType.WorkerType]dbType.WorkerType{
	mqType.WorkerType_VIDEO_ENCODER:       dbType.WorkerTypeVideoEncoder,
	mqType.WorkerType_ASCII_ENCODER:       dbType.WorkerTypeAsciiEncoder,
	mqType.WorkerType_THUMBNAIL_GENERATOR: dbType.WorkerTypeThumbnailGenerator,
	mqType.WorkerType_ASSEMBLER:           dbType.WorkerTypeAssembler,
	mqType.WorkerType_VIDEO_EXTRACTOR:     dbType.WorkerTypeVideoExtractor,
	mqType.WorkerType_AUDIO_EXTRACTOR:     dbType.WorkerTypeAudioExtractor,
	mqType.WorkerType_METADATA_EXTRACTOR:  dbType.WorkerTypeMetadataExtractor,
}
var WorkerTypeToProtoWorkerType = map[dbType.WorkerType]mqType.WorkerType{
	dbType.WorkerTypeVideoEncoder:       mqType.WorkerType_VIDEO_ENCODER,      
	dbType.WorkerTypeAsciiEncoder:       mqType.WorkerType_ASCII_ENCODER,      
	dbType.WorkerTypeThumbnailGenerator: mqType.WorkerType_THUMBNAIL_GENERATOR,
	dbType.WorkerTypeAssembler:          mqType.WorkerType_ASSEMBLER,          
	dbType.WorkerTypeVideoExtractor:     mqType.WorkerType_VIDEO_EXTRACTOR,    
	dbType.WorkerTypeAudioExtractor:     mqType.WorkerType_AUDIO_EXTRACTOR,    
	dbType.WorkerTypeMetadataExtractor:  mqType.WorkerType_METADATA_EXTRACTOR, 
}
