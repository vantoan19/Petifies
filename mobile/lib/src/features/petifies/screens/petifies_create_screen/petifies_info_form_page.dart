import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_create_screen.dart';
import 'package:mobile/src/widgets/textfield/textformfield.dart';

class PetifiesInfoFormPage extends ConsumerStatefulWidget {
  const PetifiesInfoFormPage({super.key});

  static String? titleValidator(String? title) {
    if (title == null || title.isEmpty) {
      return 'Title is required';
    }
    if (title.length > 500) {
      return 'Title cannot exceed 500 characters';
    }
    return null;
  }

  static String? petNameValidator(String? title) {
    if (title == null || title.isEmpty) {
      return 'Pet name is required';
    }
    return null;
  }

  static String? descriptionValidator(String? description) {
    if (description == null || description.isEmpty) {
      return 'Description is required';
    }
    if (description.length < 200) {
      final remainingCharacter = 200 - description.length;
      return 'Need $remainingCharacter character(s) to reach minimum length of description';
    }
    if (description.length > 5000) {
      return 'Title cannot exceed 5000 characters';
    }
    return null;
  }

  @override
  ConsumerState<PetifiesInfoFormPage> createState() =>
      _PetifiesInfoFormPageState();
}

class _PetifiesInfoFormPageState extends ConsumerState<PetifiesInfoFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _petNameController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    Future.delayed(const Duration(milliseconds: 1000), () {
      _formKey.currentState?.validate();
    });

    final title = ref.read(titleProvider);
    final description = ref.read(descriptionProvider);
    final petName = ref.read(petNameProvider);
    if (petName != null) {
      _petNameController.text = petName;
    }
    if (title != null) {
      _titleController.text = title;
    }
    if (description != null) {
      _descriptionController.text = description;
    }

    return SingleChildScrollView(
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(
              Constants.petifiesExpoloreHorizontalPadding,
              28,
              Constants.petifiesExpoloreHorizontalPadding,
              20,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  "Give us more information",
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.only(top: 8.0),
                  child: Text(
                    "Please provide a title and a description with at least 200 characters and at most 5000 characters.",
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w300,
                      color: Theme.of(context).colorScheme.secondary,
                    ),
                  ),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 4),
            child: Form(
              key: _formKey,
              child: Column(
                children: [
                  CustomTextFormField(
                    label: "Title",
                    icon: Icon(
                      Icons.title,
                      size: 18,
                    ),
                    onChange: (value) {
                      ref.read(titleProvider.notifier).state = value;
                      _formKey.currentState?.validate();
                    },
                    controller: _titleController,
                    maxLines: null,
                    validator: PetifiesInfoFormPage.titleValidator,
                  ),
                  CustomTextFormField(
                    label: "Description",
                    icon: Icon(
                      Icons.description,
                      size: 18,
                    ),
                    maxLines: null,
                    onChange: (value) {
                      ref.read(descriptionProvider.notifier).state = value;
                      _formKey.currentState?.validate();
                    },
                    controller: _descriptionController,
                    validator: PetifiesInfoFormPage.descriptionValidator,
                  ),
                  CustomTextFormField(
                    label: "Pet name",
                    icon: Icon(
                      Icons.pets,
                      size: 18,
                    ),
                    maxLines: null,
                    onChange: (value) {
                      ref.read(petNameProvider.notifier).state = value;
                      _formKey.currentState?.validate();
                    },
                    controller: _petNameController,
                    validator: PetifiesInfoFormPage.petNameValidator,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
