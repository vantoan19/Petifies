import 'dart:ffi';
import 'dart:io';

import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/proto/common/common.pb.dart';
import 'package:mobile/src/providers/service_providers.dart';
import 'package:mobile/src/services/media_service.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

final fileRepositoryProvider = Provider<IFileRepository>(
    (ref) => FileRepository(mediaService: ref.read(mediaServiceProvider)));

abstract class IFileRepository {
  Future<Either<Failure, NetworkImageModel>> uploadImage({
    required String uploaderID,
    required File image,
  });
  Future<Either<Failure, NetworkVideoModel>> uploadVideo({
    required String uploaderID,
    required File video,
  });
  Future<Either<Failure, void>> removeFile({
    required String uri,
  });
}

class FileRepository implements IFileRepository {
  final MediaService _mediaService;

  FileRepository({
    required MediaService mediaService,
  }) : _mediaService = mediaService;

  Future<Either<Failure, NetworkImageModel>> uploadImage({
    required String uploaderID,
    required File image,
  }) async {
    try {
      UploadFileResponse resp = await _mediaService.uploadImage(
        uploaderID: uploaderID,
        image: image,
      );

      NetworkImageModel imageModel = NetworkImageModel(uri: resp.uri);
      return right(imageModel);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }

  Future<Either<Failure, NetworkVideoModel>> uploadVideo({
    required String uploaderID,
    required File video,
  }) async {
    try {
      UploadFileResponse resp = await _mediaService.uploadVideo(
        uploaderID: uploaderID,
        video: video,
      );

      NetworkVideoModel videoModel = NetworkVideoModel(uri: resp.uri);
      return right(videoModel);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }

  Future<Either<Failure, void>> removeFile({
    required String uri,
  }) async {
    try {
      await _mediaService.removeFileByURI(uri: uri);
      return right(Void);
    } catch (e) {
      return left(Failure(e.toString()));
    }
  }
}
