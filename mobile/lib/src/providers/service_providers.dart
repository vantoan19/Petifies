import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/services/comment_service.dart';
import 'package:mobile/src/services/media_service.dart';
import 'package:mobile/src/services/petifies_service.dart';
import 'package:mobile/src/services/post_service.dart';
import 'package:mobile/src/services/user_service.dart';

final userServiceProvider = Provider(
  (ref) => UserService(ref: ref),
);

final postServiceProvider = Provider(
  (ref) => PostService(ref: ref),
);

final mediaServiceProvider = Provider(
  (ref) => MediaService(ref: ref),
);

final commentServiceProvider = Provider(
  (ref) => CommentService(ref: ref),
);

final petifiesServiceProvider = Provider(
  (ref) => PetifiesService(ref: ref),
);
