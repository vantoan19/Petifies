import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:mobile/src/constants/constants.dart';
import 'package:mobile/src/features/petifies/navigators/explore_navigator.dart';
import 'package:mobile/src/features/petifies/screens/petifies_details_screen.dart';
import 'package:mobile/src/providers/petifies_providers.dart';
import 'package:mobile/src/widgets/buttons/no_padding_icon_button.dart';
import 'package:mobile/src/widgets/carousel_slider/carousel_slider_with_indicator.dart';

class Petifies extends ConsumerWidget {
  const Petifies({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final petifiesInfo = ref.watch(petifiesInfoProvider);
    final width = MediaQuery.of(context).size.width -
        Constants.petifiesExpoloreHorizontalPadding * 2;

    return GestureDetector(
      behavior: HitTestBehavior.translucent,
      onTap: () => Navigator.pushNamed(
        context,
        PetifiesExploreNavigatorRoutes.petifiesDetails,
        arguments: PetifiesDetailsScreenArguments(
          petifiesData: ref.read(petifiesInfoProvider),
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(
          horizontal: Constants.petifiesExpoloreHorizontalPadding,
          vertical: 24,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            ImageCarouselSliderWithIndicators(
              images: petifiesInfo.image,
              width: width,
              height: width,
              roundCorner: true,
            ),
            Padding(
              padding: const EdgeInsets.only(top: 12),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Flexible(
                    flex: 6,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          petifiesInfo.petName,
                          style: TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.w400,
                          ),
                        ),
                        Padding(
                          padding: const EdgeInsets.only(top: 8),
                          child: Text(
                            petifiesInfo.address.toString(),
                            style: TextStyle(
                              fontSize: 14,
                              color: Theme.of(context).colorScheme.secondary,
                              fontWeight: FontWeight.w300,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  // Book mark button
                  Flexible(
                    flex: 1,
                    child: NoPaddingIconButton(
                      onPressed: () {},
                      icon: Icon(
                        FontAwesomeIcons.bookmark,
                        color: Theme.of(context).colorScheme.secondary,
                        size: 16,
                      ),
                    ),
                  ),
                  // Share button
                  Flexible(
                    flex: 1,
                    child: NoPaddingIconButton(
                      onPressed: () {},
                      icon: Icon(
                        FontAwesomeIcons.paperPlane,
                        color: Theme.of(context).colorScheme.secondary,
                        size: 16,
                      ),
                    ),
                  ),
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}
