import { useEffect, useState } from 'react';
import { api } from '../api.js';
import { setState, toast } from '../store.js';

/**
 * Anmeldung per SSO (OpenID Connect, z. B. Keycloak). Nur für Administratoren.
 *
 * Wer sich beim Anbieter anmeldet, landet im Planr-Konto mit derselben,
 * vom Anbieter bestätigten E-Mail-Adresse; gibt es keines, wird es angelegt.
 * Die Anmeldung mit Passwort bleibt daneben bestehen.
 */
export default function SSOEinstellungen() {
  const [stand, setStand] = useState(null);
  const [k, setK] = useState({ aktiv: false, aussteller: '', clientId: '', adresse: '', erlaubt: '' });
  const [geheimnis, setGeheimnis] = useState('');
  const [meldung, setMeldung] = useState(null);

  useEffect(() => {
    api('/api/sso/einstellungen')
      .then((s) => {
        setStand(s);
        setK(s.einstellung);
      })
      .catch((e) => setMeldung({ text: e.message, schlecht: true }));
  }, []);

  async function speichern(e) {
    e.preventDefault();
    try {
      const s = await api('/api/sso/einstellungen', { methode: 'PUT', daten: { ...k, geheimnis } });
      setStand(s);
      setGeheimnis('');
      setMeldung(null);
      setState({ auth: await api('/api/status') });
      toast('Gespeichert');
    } catch (err) {
      setMeldung({ text: err.message, schlecht: true });
    }
  }

  async function testen() {
    setMeldung({ text: 'Prüfe …' });
    try {
      const t = await api('/api/sso/testen', { methode: 'POST', daten: { aussteller: k.aussteller } });
      setMeldung({ text: t.meldung, schlecht: !t.ok });
    } catch (err) {
      setMeldung({ text: err.message, schlecht: true });
    }
  }

  const feld = (name) => ({ value: k[name], onChange: (e) => setK({ ...k, [name]: e.target.value }) });

  if (!stand) return <p className={meldung?.schlecht ? 'board-error' : 'board-hint'}>{meldung?.text ?? 'Wird geladen …'}</p>;

  return (
    <>
      <p className="board-hint">
        Zusätzlich zur Anmeldung mit Passwort. Konten werden über die vom Anbieter bestätigte E-Mail-Adresse
        zugeordnet. Neue Konten entstehen nur bei offener Registrierung oder gesetzter Freigabeliste.
      </p>
      {stand.notaus && <p className="board-error">Per PLANR_SSO_AUS=1 abgeschaltet.</p>}
      <form className="konto-form" onSubmit={speichern}>
        <label style={{ flexDirection: 'row', alignItems: 'center', flexBasis: '100%' }}>
          <input type="checkbox" checked={k.aktiv} onChange={(e) => setK({ ...k, aktiv: e.target.checked })} />
          „Mit SSO anmelden“ anbieten
        </label>
        <label>
          Aussteller (Issuer-URL)
          <input {...feld('aussteller')} placeholder="https://vm-debian-docker:8180/realms/homelab" />
        </label>
        <label>
          Client-ID
          <input {...feld('clientId')} placeholder="planr" />
        </label>
        <label>
          Client-Geheimnis
          <input
            type="password"
            autoComplete="new-password"
            value={geheimnis}
            onChange={(e) => setGeheimnis(e.target.value)}
            placeholder={stand.geheimnisDa ? 'gespeichert, leer lassen zum Behalten' : 'leer bei öffentlichem Client'}
          />
        </label>
        <label>
          Erlaubte Konten
          <input {...feld('erlaubt')} placeholder="leer = alle; sonst E-Mail oder Nutzername" />
        </label>
        <label>
          Adresse von Planr
          <input {...feld('adresse')} placeholder="leer = wie aufgerufen" />
        </label>
        <p className="board-hint" style={{ flexBasis: '100%' }}>
          Beim Anbieter als Weiterleitung eintragen: <code>{stand.rueckweg}</code>, nach dem Abmelden:{' '}
          <code>{stand.abmeldeweg}</code>
        </p>
        <button type="submit" className="knopf-primaer">
          Speichern
        </button>
        <button type="button" onClick={testen}>
          Verbindung testen
        </button>
      </form>
      {meldung && (
        <p role="status" className={meldung.schlecht ? 'board-error' : 'board-hint'}>
          {meldung.text}
        </p>
      )}
    </>
  );
}
