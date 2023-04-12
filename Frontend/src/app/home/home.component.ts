import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  consola: any =
    {
      entrada: '',
      salida: '',
      errores: [],
      simbolos: []
    }
  archivoSeleccionado: File|null = null;
  constructor() { }
  ngOnInit(): void {
  }

  onFileSelected(event: Event) {
    const inputFile = event.target as HTMLInputElement;
    const archivo: File = inputFile.files![0];
    if (archivo && archivo.name.endsWith('.eea')) {
      this.archivoSeleccionado = archivo;
      const lector: FileReader = new FileReader();
      lector.readAsText(this.archivoSeleccionado);
      lector.onload = (e) => {
        if (e.target && e.target.result) {
          this.mostrarContenido(e.target.result.toString());
        }
      }
    } else {
      alert('El archivo seleccionado no tiene extensión .eea');
    }
  }

  mostrarContenido(contenido: string) {
    const textarea = document.getElementById('exampleFormControlTextarea1');
    if (textarea) {
      textarea.innerHTML = contenido;
    }
  }

  compilar() {

  }
}
