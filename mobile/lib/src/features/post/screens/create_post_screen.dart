import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/feed/screens/feed_screen.dart';
import 'package:mobile/src/features/post/controllers/create_post_controller.dart';
import 'package:mobile/src/features/media/controllers/image_controller.dart';
import 'package:mobile/src/features/media/controllers/video_controller.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/uploading_post.dart';
import 'package:mobile/src/models/user_model.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/widgets/appbars/create_post_appbar.dart';
import 'package:image_picker/image_picker.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';
import 'package:mobile/src/widgets/images/image_upload_viewer.dart';
import 'package:mobile/src/widgets/pickers/visibility_picker.dart';
import 'package:uuid/uuid.dart';
import 'package:video_player/video_player.dart';

class CreatePostScreen extends ConsumerStatefulWidget {
  const CreatePostScreen({super.key});

  @override
  ConsumerState<ConsumerStatefulWidget> createState() =>
      _CreatePostScreenState();
}

class _CreatePostScreenState extends ConsumerState<CreatePostScreen> {
  TextEditingController _textController = TextEditingController();
  List<File> _images = [];
  List<VideoPlayerController> _videoControllers = [];
  String _visibility = "public";
  List<Future<Either<Failure, NetworkImageModel>>> _imagesFutures = [];
  List<Future<Either<Failure, NetworkVideoModel>>> _videosFutures = [];

  final ImagePicker _picker = ImagePicker();
  late FocusNode _textContentFocusNode;

  Future<void> _pickImages(String userID) async {
    if (_images.length + _videoControllers.length >= 4) {
      return;
    }

    final List<XFile>? images = await _picker.pickMultiImage();
    if (images != null) {
      setState(() {
        _images.addAll(images.map((img) => File(img.path)));
        _imagesFutures.addAll(
          images.map(
            (img) {
              return ref
                  .read(imageControllerProvider.notifier)
                  .uploadImage(uploaderID: userID, image: File(img.path));
            },
          ),
        );
      });
    }
    _textContentFocusNode.requestFocus();
  }

  Future<void> _pickVideo(String userID) async {
    if (_images.length + _videoControllers.length >= 4) {
      return;
    }

    final XFile? video = await _picker.pickVideo(source: ImageSource.gallery);
    if (video != null) {
      setState(() {
        _videoControllers.add(
          VideoPlayerController.file(File(video.path))
            ..initialize().then((_) {
              setState(() {});
            }),
        );

        _videosFutures.add(ref
            .read(videoControllerProvider.notifier)
            .uploadVideo(uploaderID: userID, video: File(video.path)));
      });
    }
    _textContentFocusNode.requestFocus();
  }

  Future<void> _takePicture(String userID) async {
    if (_images.length + _videoControllers.length >= 4) {
      return;
    }

    final XFile? image = await _picker.pickImage(source: ImageSource.camera);
    if (image != null) {
      setState(() {
        _images.add(File(image.path));

        _imagesFutures.add(ref
            .read(imageControllerProvider.notifier)
            .uploadImage(uploaderID: userID, image: File(image.path)));
      });
    }
    _textContentFocusNode.requestFocus();
  }

  Future<void> _takeVideo(String userID) async {
    if (_images.length + _videoControllers.length >= 4) {
      return;
    }

    final XFile? video = await _picker.pickVideo(source: ImageSource.camera);
    if (video != null) {
      setState(() {
        _videoControllers.add(
          VideoPlayerController.file(File(video.path))
            ..initialize().then((_) {
              setState(() {});
            }),
        );

        _videosFutures.add(ref
            .read(videoControllerProvider.notifier)
            .uploadVideo(uploaderID: userID, video: File(video.path)));
      });
    }
    _textContentFocusNode.requestFocus();
  }

  void _removeImage(int idx) async {
    if (idx >= _images.length) {
      return;
    }

    Future<Either<Failure, NetworkImageModel>> imgFuture = _imagesFutures[idx];
    setState(() {
      _images.removeAt(idx);
      _imagesFutures.removeAt(idx);
    });
    (await imgFuture).fold((l) => null, (r) {
      ref.read(imageControllerProvider.notifier).removeImage(uri: r.uri);
    });
  }

  void _submitPostCreationRequest(BasicUserInfoModel author) {
    if (_images.length + _videoControllers.length > 4) {
      return;
    }
    if (_textController.text == "" &&
        _images.length == 0 &&
        _videoControllers.length == 0) {
      return;
    }

    final tempId = Uuid().v4();
    ref
        .read(newlyCreatedPostsProvider.notifier)
        .addNewlyCreatedPost(UploadingPostModel(
          tempId: tempId,
          owner: author,
          postActivity: "post",
          textContent: _textController.text,
          images: _images,
          videos: _videoControllers,
          createdAt: DateTime.now(),
        ));

    ref.read(createPostControllerProvider.notifier).createPost(
          tempId: tempId,
          author: author,
          visibility: _visibility,
          activity: "post",
          textContent: _textController.text,
          imageFutures: _imagesFutures,
          videoFutures: _videosFutures,
        );

    NavigatorUtil.goBack(context);
  }

  @override
  void initState() {
    super.initState();
    _textContentFocusNode = FocusNode();
  }

  @override
  void dispose() {
    // Clean up the focus node when the Form is disposed.
    _textContentFocusNode.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final user = ref.watch(myUserProvider);

    final err = user.maybeWhen(
      error: (e, stackTrace) => e.toString(),
      orElse: () => "no err",
    );

    final userInfo = user.when(
      data: (data) => data,
      error: (e, stackTrace) => null,
      loading: () => null,
    );

    if (userInfo == null) {
      return Placeholder();
    }

    return Scaffold(
      appBar: CreatePostAppBar(
        addPostAction: () => _submitPostCreationRequest(BasicUserInfoModel(
          id: userInfo.id,
          email: userInfo.email,
          firstName: userInfo.firstName,
          lastName: userInfo.lastName,
        )),
      ),
      body: Padding(
        padding: EdgeInsets.symmetric(horizontal: 16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.max,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 0, 0, 16),
              child: Row(
                children: [
                  // User Avatar
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 0, 12, 0),
                    child: (userInfo.userAvatar != null)
                        ? CircleAvatar(
                            backgroundImage: NetworkImage(
                              userInfo.userAvatar!,
                            ),
                            backgroundColor: Colors.transparent,
                            radius: 24,
                          )
                        : CircleAvatar(
                            backgroundImage: AssetImage(
                              Constants.defaultAvatarPng,
                            ),
                            backgroundColor: Colors.transparent,
                            radius: 24,
                          ),
                  ),
                  // Visibility selection
                  VisibilityPicker(
                    initialValue: _visibility,
                    onVisibilityChangedFunc: (value) {
                      if (value != null) {
                        _visibility = value;
                      }
                    },
                  ),
                ],
              ),
            ),
            // Post Content
            Expanded(
              child: SingleChildScrollView(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    TextFormField(
                      autofocus: true,
                      controller: _textController,
                      decoration: InputDecoration(
                          hintText: 'What\'s happening?',
                          border: OutlineInputBorder(
                            borderSide: BorderSide.none,
                          ),
                          contentPadding: EdgeInsets.symmetric(vertical: 10)),
                      style: TextStyle(
                        fontSize: 20,
                      ),
                      maxLines: null,
                      focusNode: _textContentFocusNode,
                      keyboardType: TextInputType.multiline,
                    ),
                    SizedBox(
                      height: MediaQuery.of(context).size.width * 0.5,
                      child: ListView(
                        scrollDirection: Axis.horizontal,
                        children: [
                          ..._images.asMap().entries.map((entry) {
                            return ImageUploadViewer(
                              image: entry.value,
                              uploaderID: userInfo.id,
                              removeAction: () => _removeImage(entry.key),
                            );
                          }).toList(),
                          ..._videoControllers.map(
                            (video) => AspectRatio(
                              aspectRatio: video.value.aspectRatio,
                              child: VideoPlayer(video),
                            ),
                          )
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Row(
                  children: [
                    IconButton(
                      icon: Icon(Icons.camera_alt_outlined),
                      onPressed: () => _takePicture(userInfo.id),
                    ),
                    IconButton(
                      icon: Icon(Icons.videocam_outlined),
                      onPressed: () => _takeVideo(userInfo.id),
                    ),
                    IconButton(
                      icon: Icon(Icons.photo_library_outlined),
                      onPressed: () => _pickImages(userInfo.id),
                    ),
                    IconButton(
                      icon: Icon(Icons.video_library_outlined),
                      onPressed: () => _pickVideo(userInfo.id),
                    ),
                  ],
                ),
                IconButton(
                  icon: Icon(
                    Icons.keyboard_outlined,
                  ),
                  onPressed: () {
                    if (_textContentFocusNode.hasFocus) {
                      _textContentFocusNode.unfocus();
                    } else {
                      _textContentFocusNode.requestFocus();
                    }
                  },
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
