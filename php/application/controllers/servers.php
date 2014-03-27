<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');

class Servers extends CI_Controller {
	public function listservers()
	{
		//$list = json_decode(file_get_contents(GO_PATH.'/servers/list'));

		$list = array(
			array(
				'name' => 'server1'
				,'ip' => '1.2.3.4'
				,'status' => 'running'
				,'id' => 1
			)
			,array(
				'name' => 'server2'
				,'ip' => '1.2.3.4'
				,'status' => 'down'
				,'id' => 2
			)
			,array(
				'name' => 'server3'
				,'ip' => '1.2.3.4'
				,'status' => 'full'
				,'id' => 3
			)
			,array(
				'name' => 'server4'
				,'ip' => '1.2.3.4'
				,'status' => 'high'
				,'id' => 4
			)
			,array(
				'name' => 'server5'
				,'ip' => '1.2.3.4'
				,'status' => 'drunk'
				,'id' => 5
			)
		);

		//$this->load->view('servers/listservers',array('servers' => $list));
		$this->layout->setLayout('layouts/main');
		$this->layout->view('servers/listservers',array(
			'js_file' => 'AdminLTE/dashboard'
			,'servers' => $list
		));
	}

	function addserver(){
		$this->layout->setLayout('layouts/main');
		$this->layout->view('servers/addserver',array(
			'js_file' => ''
		));
	}

	function doaddserver(){

	}
}

/* End of file welcome.php */
/* Location: ./application/controllers/welcome.php */