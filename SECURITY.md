# Security Policy

## Podporne različice

| Različica | Podprta |
|-----------|---------|
| 0.13.x    | ✅      |

## Prijava ranljivosti

Če odkrijete varnostno ranljivost, prosimo, da je **NE objavljate javno**.

Pošljite opis na:
- GitHub Issues (z oznako `security`)
- ali direktno lastniku repozitorija

Vključite:
- Opis ranljivosti
- Korake za repliciranje
- Morebitne predloge za popravek

Odgovorili bomo v roku **48 ur**.

## Varnostne mere v kodi

- JWT secret se generira naključno ob zagonu (če ni nastavljen)
- PIN kod so SHA-256 soljeno hashirani
- Passwordi so bcrypt hashirani
- Rate limiting na auth endpointih
- Uporabniški vnos ekraniran v MongoDB queryjih
- CORS konfiguriran
- Docker teče kot non-root uporabnik
