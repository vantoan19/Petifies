import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/screens/petifies_create_screen/petifies_create_screen.dart';
import 'package:mobile/src/models/petifies.dart';

class PetifiesTypeSelectionPage extends ConsumerWidget {
  static final types = [
    PetifiesType.DOG_WALKING,
    PetifiesType.CAT_PLAYING,
    PetifiesType.DOG_SITTING,
    PetifiesType.CAT_SITTING,
    PetifiesType.DOG_ADOPTION,
    PetifiesType.CAT_ADOPTION,
  ];

  const PetifiesTypeSelectionPage({
    Key? key,
  }) : super(key: key);

  String getLabel(PetifiesType type) {
    switch (type) {
      case PetifiesType.DOG_WALKING:
        return "Dog walking";
      case PetifiesType.CAT_PLAYING:
        return "Cat playing";
      case PetifiesType.DOG_SITTING:
        return "Dog sitting";
      case PetifiesType.CAT_SITTING:
        return "Cat sitting";
      case PetifiesType.DOG_ADOPTION:
        return "Dog adoption";
      case PetifiesType.CAT_ADOPTION:
        return "Cat adoption";
      case PetifiesType.UNKNOWN:
        return "Unknown";
    }
  }

  IconData getIcon(PetifiesType type) {
    switch (type) {
      case PetifiesType.DOG_WALKING:
        return FontAwesomeIcons.dog;
      case PetifiesType.CAT_PLAYING:
        return FontAwesomeIcons.cat;
      case PetifiesType.DOG_SITTING:
        return FontAwesomeIcons.shieldDog;
      case PetifiesType.CAT_SITTING:
        return FontAwesomeIcons.shieldCat;
      case PetifiesType.DOG_ADOPTION:
        return FontAwesomeIcons.bone;
      case PetifiesType.CAT_ADOPTION:
        return FontAwesomeIcons.bowlFood;
      case PetifiesType.UNKNOWN:
        return FontAwesomeIcons.question;
    }
  }

  void _changePetifiesType(PetifiesType type, WidgetRef ref) {
    if (ref.read(petifiesTypeProvider) == type) {
      ref.read(petifiesTypeProvider.notifier).state = null;
    } else {
      ref.read(petifiesTypeProvider.notifier).state = type;
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.fromLTRB(
            Constants.petifiesExpoloreHorizontalPadding,
            28,
            Constants.petifiesExpoloreHorizontalPadding,
            0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Padding(
              padding: const EdgeInsets.only(bottom: 20),
              child: Text(
                "Which one of these best describes your Petifies?",
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.w500,
                  letterSpacing: 0.5,
                ),
              ),
            ),
            Expanded(
              child: GridView.count(
                crossAxisCount: 2,
                mainAxisSpacing: 12,
                crossAxisSpacing: 12,
                childAspectRatio: 1.4,
                children: [
                  ...types
                      .map(
                        (e) => _PetifiesTypeSelectionButton(
                          action: () => _changePetifiesType(e, ref),
                          label: getLabel(e),
                          icon: getIcon(e),
                          type: e,
                        ),
                      )
                      .toList(),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _PetifiesTypeSelectionButton extends ConsumerWidget {
  final IconData icon;
  final String label;
  final PetifiesType type;
  final VoidCallback action;
  const _PetifiesTypeSelectionButton({
    Key? key,
    required this.icon,
    required this.label,
    required this.type,
    required this.action,
  }) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isSelected = ref.watch(petifiesTypeProvider) == type;

    return GestureDetector(
      onTap: action,
      child: Card(
        child: Padding(
          padding: const EdgeInsets.only(left: 16),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Icon(icon),
              Text(
                label,
                style: TextStyle(
                  fontSize: 16,
                  color: Theme.of(context).colorScheme.inversePrimary,
                  letterSpacing: 0.5,
                ),
              )
            ],
          ),
        ),
        shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
            side: BorderSide(
              width: 2,
              color: Theme.of(context)
                  .colorScheme
                  .secondary
                  .withOpacity(isSelected ? 1 : 0.25),
            )),
        color: isSelected
            ? Theme.of(context).colorScheme.secondary.withOpacity(0.1)
            : Theme.of(context).scaffoldBackgroundColor,
      ),
    );
  }
}
