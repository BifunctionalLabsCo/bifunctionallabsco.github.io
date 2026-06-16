import rss from '@astrojs/rss';

// Manually maintained list -- add an entry here whenever a new expedition post is published.
const expeditions = [
  {
    title: 'Andaman Sands',
    description: "Phuket sunsets, rock pools, and island hopping to Ko Phi Phi on Thailand's Andaman coast.",
    link: '/expeditions/andaman/',
    pubDate: new Date('2019-01-07'),
  },
  {
    title: 'Bangkok Airwaves',
    description: "Wat Pho, Wat Arun, and the temples of Bangkok's Phra Nakhon district.",
    link: '/expeditions/bangkok/',
    pubDate: new Date('2019-03-24'),
  },
  {
    title: 'Verona: Of Bridges & Brawls',
    description: "The Arena di Verona, Piazza Bra, and the Roman ruins of Verona's old city.",
    link: '/expeditions/verona/',
    pubDate: new Date('2020-10-08'),
  },
  {
    title: 'Into the Alps: Innsbruck',
    description: 'A weekend in Innsbruck: Altstadt, the City Tower, the Nordkette trail, and Castle Ambras.',
    link: '/expeditions/innsbruck/',
    pubDate: new Date('2021-10-09'),
  },
  {
    title: 'Into the Alps: Salzburg',
    description: 'A weekend in Salzburg: Untersberg, Hohensalzburg fortress, and the Sound of Music filming locations.',
    link: '/expeditions/salzburg/',
    pubDate: new Date('2021-10-09'),
  },
];

export async function GET(context) {
  return rss({
    title: 'Bifunctional Expeditions',
    description: 'Travel stories from Saumya and Zubin.',
    site: context.site,
    items: expeditions,
  });
}
