import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/models/comment.dart';
import 'package:mobile/src/models/uploading_comment.dart';

final commentInfoProvider =
    Provider<CommentModel>((ref) => throw UnimplementedError());

final uploadingCommentInfoProvider =
    Provider<UploadingCommentModel>((ref) => throw UnimplementedError());

// ancestor comments, include itself
final ancestorCommentsProvider =
    Provider<List<CommentModel>>((ref) => throw UnimplementedError());
