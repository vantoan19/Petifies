// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:image_picker/image_picker.dart';
import 'package:mobile/src/exceptions/failure.dart';
import 'package:mobile/src/features/comment/controller/create_comment_controller.dart';
import 'package:mobile/src/features/comment/controller/list_comments_controller.dart';
import 'package:mobile/src/features/media/controllers/image_controller.dart';
import 'package:mobile/src/models/basic_user_info.dart';
import 'package:mobile/src/models/image.dart';
import 'package:mobile/src/models/uploading_comment.dart';
import 'package:mobile/src/models/video.dart';
import 'package:mobile/src/providers/user_model_providers.dart';
import 'package:mobile/src/utils/navigation.dart';
import 'package:mobile/src/widgets/appbars/create_comment_appbar.dart';
import 'package:mobile/src/widgets/buttons/love_react_button.dart';
import 'package:mobile/src/widgets/comment/comment.dart';
import 'package:mobile/src/widgets/images/image_upload_viewer.dart';
import 'package:mobile/src/widgets/posts/post.dart';
import 'package:mobile/src/widgets/user_avatar/user_avatar.dart';
import 'package:uuid/uuid.dart';
import 'package:video_player/video_player.dart';

class CreateCommentScreen extends ConsumerStatefulWidget {
  const CreateCommentScreen({
    Key? key,
  }) : super(key: key);

  @override
  ConsumerState<ConsumerStatefulWidget> createState() =>
      _CreateCommentScreenState();
}

class _CreateCommentScreenState extends ConsumerState<CreateCommentScreen> {
  TextEditingController _textController = TextEditingController();
  File? _image = null;
  VideoPlayerController? _videoController = null;
  Future<Either<Failure, NetworkImageModel>>? _imageFuture = null;
  Future<Either<Failure, NetworkVideoModel>>? _videoFuture = null;

  final ImagePicker _picker = ImagePicker();
  late FocusNode _textContentFocusNode;

  Future<void> _pickImage(String userID) async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      setState(() {
        _image = File(image.path);
        _imageFuture = ref
            .read(imageControllerProvider.notifier)
            .uploadImage(uploaderID: userID, image: File(image.path));
      });
    }
    _textContentFocusNode.requestFocus();
  }

  Future<void> _takePicture(String userID) async {
    final XFile? image = await _picker.pickImage(source: ImageSource.camera);
    if (image != null) {
      setState(() {
        _image = File(image.path);
        _imageFuture = ref
            .read(imageControllerProvider.notifier)
            .uploadImage(uploaderID: userID, image: File(image.path));
      });
    }
    _textContentFocusNode.requestFocus();
  }

  void _removeImage() async {
    Future<Either<Failure, NetworkImageModel>>? imgFuture = _imageFuture;
    setState(() {
      _image = null;
      _imageFuture = null;
    });
    if (imgFuture != null) {
      (await imgFuture).fold((l) => null, (r) {
        ref.read(imageControllerProvider.notifier).removeImage(uri: r.uri);
      });
    }
  }

  void _submitCommentCreationRequest(BasicUserInfoModel author) {
    if (_image != null && _videoController != null) {
      return;
    }
    if (_textController.text == "" &&
        _image == null &&
        _videoController == null) {
      return;
    }

    final isPostTarget = ref.read(isPostContextProvider);
    final targetID = isPostTarget
        ? ref.read(postInfoProvider).id
        : ref.read(commentInfoProvider).id;

    final tempId = Uuid().v4();
    ref
        .read(newlyCreatedCommentsProvider(parentID: targetID).notifier)
        .addNewlyCreatedComment(UploadingCommentModel(
          tempId: tempId,
          owner: author,
          textContent: _textController.text,
          image: _image,
          video: _videoController,
          createdAt: DateTime.now(),
        ));
    ref.read(createCommentControllerProvider.notifier).createComment(
          tempId: tempId,
          author: author,
          postID: ref.read(postInfoProvider).id,
          parentID: targetID,
          isPostParent: isPostTarget,
          textContent: _textController.text,
          imageFuture: _imageFuture,
          videoFuture: _videoFuture,
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

    final isPostTarget = ref.read(isPostContextProvider);
    final targetOwner = isPostTarget
        ? ref.read(postInfoProvider).owner
        : ref.read(commentInfoProvider).owner;

    return Scaffold(
      appBar: CreateCommentAppBar(
        addPostAction: () => _submitCommentCreationRequest(
          BasicUserInfoModel(
            id: userInfo.id,
            email: userInfo.email,
            firstName: userInfo.firstName,
            lastName: userInfo.lastName,
          ),
        ),
        targetOwner: targetOwner.firstName + " " + targetOwner.lastName,
        targetType: isPostTarget ? "post" : "comment",
      ),
      body: Column(
        children: [
          Expanded(
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _CreateCommentUserAvatar(
                  userAvatar: userInfo.userAvatar,
                ),
                Expanded(
                  child: SingleChildScrollView(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _CreateCommentTextFormField(
                          textController: _textController,
                          textContentFocusNode: _textContentFocusNode,
                        ),
                        if (_image != null)
                          ImageUploadViewer(
                            image: _image!,
                            uploaderID: userInfo.id,
                            removeAction: () => _removeImage(),
                          )
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
          _ControlButtons(
            takePictureCallback: () => _takePicture(userInfo.id),
            choosePictureCallBack: () => _pickImage(userInfo.id),
            disableTakePicture: (_image == null ? false : true),
            disableChoosePicture: (_image == null ? false : true),
            textContentFocusNode: _textContentFocusNode,
          )
        ],
      ),
    );
  }
}

class _CreateCommentUserAvatar extends UserAvatar {
  final String? userAvatar;

  const _CreateCommentUserAvatar({
    Key? key,
    this.userAvatar = null,
  }) : super(
          key: key,
          userAvatar: userAvatar,
          padding: const EdgeInsets.fromLTRB(16, 0, 16, 0),
        );
}

class _CreateCommentTextFormField extends StatelessWidget {
  final TextEditingController textController;
  final FocusNode textContentFocusNode;

  const _CreateCommentTextFormField({
    Key? key,
    required this.textController,
    required this.textContentFocusNode,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      autofocus: true,
      controller: textController,
      decoration: InputDecoration(
          hintText: 'Write your reply...',
          border: OutlineInputBorder(
            borderSide: BorderSide.none,
          ),
          contentPadding: EdgeInsets.symmetric(vertical: 10)),
      style: TextStyle(
        fontSize: 20,
      ),
      maxLines: null,
      focusNode: textContentFocusNode,
      keyboardType: TextInputType.multiline,
    );
  }
}

class _ControlButtons extends StatelessWidget {
  final VoidCallback takePictureCallback;
  final VoidCallback choosePictureCallBack;
  final bool disableTakePicture;
  final bool disableChoosePicture;
  final FocusNode textContentFocusNode;

  const _ControlButtons({
    Key? key,
    required this.takePictureCallback,
    required this.choosePictureCallBack,
    required this.disableTakePicture,
    required this.disableChoosePicture,
    required this.textContentFocusNode,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Row(
          children: [
            IconButton(
              icon: Icon(Icons.camera_alt_outlined),
              onPressed: takePictureCallback,
            ),
            IconButton(
              icon: Icon(Icons.photo_library_outlined),
              onPressed: choosePictureCallBack,
            ),
          ],
        ),
        IconButton(
          icon: Icon(
            Icons.keyboard_outlined,
          ),
          onPressed: () {
            if (textContentFocusNode.hasFocus) {
              textContentFocusNode.unfocus();
            } else {
              textContentFocusNode.requestFocus();
            }
          },
        ),
      ],
    );
  }
}
