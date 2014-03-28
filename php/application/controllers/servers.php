<?php if ( ! defined('BASEPATH')) exit('No direct script access allowed');

class Servers extends CI_Controller {
	public function listservers()
	{
		$post = curl_init();

   curl_setopt($post, CURLOPT_URL, trim(GO_PATH.'server/list') );
   curl_setopt($post, CURLOPT_RETURNTRANSFER, 1);
   curl_setopt($post, CURLOPT_HTTPHEADER, array('USER_ID: '.$this->session->userdata('user_id'), 'TOKEN: '.$this->session->userdata('email')));

   $result = curl_exec($post);

   $header = curl_getinfo ( $post );

   //echo '<pre>'; print_r($header);echo '</pre>';

   curl_close($post);
		//$list = json_decode(file_get_contents(GO_PATH.'server/list'));
   		$list = json_decode($result);
   		//echo '<pre>'; print_r($list1);echo '</pre>'; exit;

		/*$list = array(
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
		);*/

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

		$post = $this->input->post();

		//lets load user model
		//u is alias of User_db
		$this->load->model('Server_db','u');

		//lets create new server by inserting into db
		$server_id = $this->u->insertNewServer($post['serverName'],$post['ip']);

		if($server_id <= 0){
			echo json_encode(array(
				'status' => 'error'
				,'msg'	=> $server_id == 0 ? 'Somer error occured. Please try after some time.' : 'A server already exists with this name. Please user correct name or ip to add server.'
			));
			exit;
		}

		$random_string = substr(sha1(rand()), 0, 15);


		//LETS REGISTER IN GO SERVER
		$status = post_to_go('add','server',array("input"=>json_encode(array(
			'name' => $post['serverName']
			,'ip' => $post['ip']
			,'key' => $random_string
		))),array('USER_ID: '.$this->session->userdata('user_id'), 'TOKEN: '.$this->session->userdata('email')));

		echo json_encode(array(
			'status' => 'success'
			,'message' => 'You are successfully registered.'
			,'server_id' => $server_id
			,'resp' => $status
		));
		exit;
	}

	function getserverstatus(){
			/*$post = curl_init();

   curl_setopt($post, CURLOPT_URL, trim(GO_PATH.'server/status') );
   curl_setopt($post, CURLOPT_RETURNTRANSFER, 1);
   curl_setopt($post, CURLOPT_HTTPHEADER, array('USER_ID: '.$this->session->userdata('user_id'), 'TOKEN: '.$this->session->userdata('email')));

   $result = curl_exec($post);

   $header = curl_getinfo ( $post );

   curl_close($post);
   		*/

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
		$this->layout->view('servers/getserverstatus',array(
			'js_file' => 'AdminLTE/dashboard'
			,'servers' => $list
		));
}
}

/* End of file welcome.php */
/* Location: ./application/controllers/welcome.php */
